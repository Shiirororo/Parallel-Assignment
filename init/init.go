package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitRedis(config Config) *redis.Client {
	r := config.Redis
	rdb := redis.NewClient(&redis.Options{
		Password: r.Password,
		DB:       r.Database,
		PoolSize: 10,
	})
	fmt.Println("Connected to redis")
	return rdb
}

func InitDB(config Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(
			config.Mongo.URI,
		),
	)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB")
	return client

}
