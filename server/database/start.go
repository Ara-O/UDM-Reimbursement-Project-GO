package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisDB *redis.Client

func startRedisDatabase() error {
	url := os.Getenv("REDIS_URL")

	opts, err := redis.ParseURL(url)
	if err != nil {
		return err
	}

	redisDB = redis.NewClient(opts)

	fmt.Println("Redis server successfully started")
	return nil
}

func StartDatabase() error {
	if err := startRedisDatabase(); err != nil {
		return err
	}

	return nil
}

func GetRedisDatabaseConnection() *redis.Client {
	return redisDB
}
