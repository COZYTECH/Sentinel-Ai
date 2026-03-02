package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if _, err := RDB.Ping(Ctx).Result(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	log.Println("Connected to Redis")
}