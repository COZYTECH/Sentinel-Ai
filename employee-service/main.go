package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/COZYTECH/Sentinel-Ai/shared/events"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if _, err := RDB.Ping(Ctx).Result(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	log.Println("Connected to Redis")
}

// Simulate employee actions every few seconds
func simulateEmployeeActions() {
	for {
		action := events.EmployeeAccessEvent{
			EmployeeID: "EMP001",
			UserID:     "USER123",
			Action:     "VIEW_ACCOUNT",
			Timestamp:  time.Now().Format(time.RFC3339),
		}
		data, _ := json.Marshal(action)
		if err := RDB.XAdd(Ctx, &redis.XAddArgs{
			Stream: "employee-activity-stream",
			Values: map[string]interface{}{"data": data},
		}).Err(); err != nil {
			log.Println("Failed to push employee event:", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	InitRedis()
	log.Println("Employee Service running...")
	simulateEmployeeActions()
}