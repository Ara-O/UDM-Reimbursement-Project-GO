package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx = context.Background()
var redisDB *redis.Client
var mongoDB *mongo.Client

func startMongoDatabase() error {
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URL"))

	mongoDB, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	defer func() {
		if err = mongoDB.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if err = mongoDB.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	fmt.Println("MongoDB server successfully started!")
	return nil
}

func startRedisDatabase() error {
	url := os.Getenv("REDIS_URL")

	opts, err := redis.ParseURL(url)
	if err != nil {
		return err
	}

	redisDB = redis.NewClient(opts)

	fmt.Println("Redis server successfully started!")
	return nil
}

func StartDatabase() error {
	if err := startRedisDatabase(); err != nil {
		return err
	}

	if err := startMongoDatabase(); err != nil {
		return err
	}

	return nil
}

func GetRedisDatabaseConnection() *redis.Client {
	return redisDB
}

func Close() {
	fmt.Println("close")
}
