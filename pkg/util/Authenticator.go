package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bakare-dev/simple-bank-api/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte(config.Settings.Security.JWTSecret)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func IsPasswordHashed(password string) bool {
	return strings.HasPrefix(password, "$2a$") || strings.HasPrefix(password, "$2b$") || strings.HasPrefix(password, "$2y$")
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWTToken(userID string, role string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(5 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	expiresAt := time.Now().Add(5 * time.Hour).Format(time.RFC3339)

	tokenData := map[string]interface{}{
		"accessToken": signedToken,
		"expiresAt":   expiresAt,
	}

	tokenDataJSON, err := json.Marshal(tokenData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token data: %w", err)
	}

	err = client.SetEx(context.Background(), fmt.Sprintf("token:%s", signedToken), tokenDataJSON, time.Duration(18000)*time.Second).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store token in Redis: %w", err)
	}

	return map[string]interface{}{
		"accessToken": signedToken,
		"expiresAt":   expiresAt,
	}, nil
}

func VerifyJWTToken(tokenString string) (map[string]interface{}, error) {
	cachedToken, err := client.Get(context.Background(), fmt.Sprintf("token:%s", tokenString)).Result()
	if err == redis.Nil {
		return nil, errors.New("token invalid or expired")
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve token from Redis: %w", err)
	}

	var tokenData map[string]interface{}
	if err := json.Unmarshal([]byte(cachedToken), &tokenData); err != nil {
		return nil, fmt.Errorf("failed to parse cached token data: %w", err)
	}

	expiresAtStr, ok := tokenData["expiresAt"].(string)
	if !ok {
		return nil, errors.New("invalid token data: missing expiration time")
	}

	expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration time: %w", err)
	}

	if time.Now().After(expiresAt) {
		return nil, errors.New("token expired")
	}

	accessToken, ok := tokenData["accessToken"].(string)
	if !ok {
		return nil, errors.New("invalid token data: missing access token")
	}

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalid or expired")
	}

	return claims, nil
}

func InvalidateToken(tokenString string) error {
	err := client.Del(context.Background(), fmt.Sprintf("token:%s", tokenString)).Err()
	if err != nil {
		return fmt.Errorf("logout failed")
	}
	return nil
}
