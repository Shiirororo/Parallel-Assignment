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
		Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
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
	uri := fmt.Sprintf("mongodb://%s:%d", config.Mongo.Host, config.Mongo.Port)
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
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
