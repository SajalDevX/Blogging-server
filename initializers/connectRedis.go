package initializers

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	ctx := context.Background()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}
