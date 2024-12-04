package util

import (
	"context"
	"fmt"
	"log"

	"github.com/bakare-dev/simple-bank-api/pkg/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	var ctx = context.Background()
	redisConfig := config.Settings.Infrastructure.Redis

	client := redis.NewClient(&redis.Options{
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
