package util

import (
	"errors"
	"time"

	"github.com/bakare-dev/simple-bank-api/config"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	cost := config.Settings.Security.BcryptCost

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GeneratePASETOToken(userID string) (string, error) {
	privateKey := []byte(config.Settings.Security.PasetoSecret)

	expirationTime := time.Now().Add(time.Hour * 5)

	token := paseto.NewV2()

	claims := map[string]interface{}{
		"sub": userID,
		"exp": expirationTime.Unix(),
	}

	encodedToken, err := token.Encrypt(privateKey, claims, nil)

	if err != nil {
		return "", err
	}

	return encodedToken, nil
}

func VerifyPASETOToken(tokenString string) (string, error) {
	privateKey := []byte(config.Settings.Security.PasetoSecret)

	token := paseto.NewV2()

	var claims map[string]interface{}

	err := token.Decrypt(tokenString, privateKey, &claims, nil)

	if err != nil {
		return "", err
	}

	expirationTime, ok := claims["exp"].(float64)

	if !ok || time.Unix(int64(expirationTime), 0).Before(time.Now()) {
		return "", errors.New("token has expired")
	}

	userID, ok := claims["sub"].(string)

	if !ok {
		return "", errors.New("invalid token claims")
	}

	return userID, nil
}
