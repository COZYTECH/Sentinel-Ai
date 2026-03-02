package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/COZYTECH/Sentinel-Ai/shared/events"
	"github.com/gin-gonic/gin"
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

func publishTransaction(tx events.TransactionCreatedEvent) error {
	data, _ := json.Marshal(tx)
	return RDB.XAdd(Ctx, &redis.XAddArgs{
		Stream: "transaction-stream",
		Values: map[string]interface{}{"data": data},
	}).Err()
}

func publishEmployeeAction(e events.EmployeeAccessEvent) error {
	data, _ := json.Marshal(e)
	return RDB.XAdd(Ctx, &redis.XAddArgs{
		Stream: "employee-activity-stream",
		Values: map[string]interface{}{"data": data},
	}).Err()
}

func main() {
	InitRedis()
	r := gin.Default()

	// POST /transactions
	r.POST("/transactions", func(c *gin.Context) {
		var tx events.TransactionCreatedEvent
		if err := c.ShouldBindJSON(&tx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tx.CreatedAt = time.Now().Format(time.RFC3339)
		if err := publishTransaction(tx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish transaction"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "transaction published"})
	})

	// POST /employee-actions
	r.POST("/employee-actions", func(c *gin.Context) {
		var e events.EmployeeAccessEvent
		if err := c.ShouldBindJSON(&e); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		e.Timestamp = time.Now().Format(time.RFC3339)
		if err := publishEmployeeAction(e); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish employee action"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "employee action published"})
	})

	r.Run(":8080")
	log.Println("API Gateway running on :8080")
}