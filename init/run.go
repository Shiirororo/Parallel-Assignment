package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run() (*redis.Client, *mongo.Client) {
	config := LoadConfig()

	rdb := InitRedis(config)
	db := InitDB(config)
	return rdb, db
}
