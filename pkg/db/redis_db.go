package db

import (
	"fmt"
	"github.com/dionofrizal88/dsi/auth/config"
	"log"

	"github.com/go-redis/redis/v8"
)

// NewRedisConnection will create a redis DB connection.
func NewRedisConnection(conf config.Configuration) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort),
		DB:   conf.RedisDB,
	})

	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)

		return nil, err
	}

	return redisClient, nil
}
