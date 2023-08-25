package services

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func OpenRedis(url string) *redis.Client {
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return client
}
