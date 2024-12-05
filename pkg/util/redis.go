package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bakare-dev/simple-bank-api/pkg/config"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func ConnectRedis() *redis.Client {
	var ctx = context.Background()
	redisConfig := config.Settings.Infrastructure.Redis

	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return client
}

func SetKey(ctx context.Context, key string, value interface{}, expirationInSec int64) error {
	cmd := client.Set(ctx, key, value, 0)
	if expirationInSec > 0 {
		cmd = client.SetEx(ctx, key, value, time.Duration(expirationInSec)*time.Second)
	}
	_, err := cmd.Result()
	if err != nil {
		return err
	}
	return nil
}

func GetKey(ctx context.Context, key string) (string, error) {
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func DelKey(ctx context.Context, key string) error {
	_, err := client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
