package main

import (
	"log"
	"net/http"
	"time"

	"encoding/json"

	"github.com/COZYTECH/Sentinel-Ai/shared/events"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CreateTransactionRequest struct {
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Country   string  `json:"country"`
	DeviceID  string  `json:"device_id"`
	IPAddress string  `json:"ip_address"`
}

func CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactionID := uuid.New().String()

	// TODO:
	// 1. Save to MySQL
	// 2. Publish event to Redis


	query := `
	INSERT INTO transactions 
	(id, user_id, amount, currency, country, device_id, ip_address, status)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(
		query,
		transactionID,
		req.UserID,
		req.Amount,
		req.Currency,
		req.Country,
		req.DeviceID,
		req.IPAddress,
		"PENDING",
	)



		event := events.TransactionCreatedEvent{
		Event:         "transaction.created",
		TransactionID: transactionID,
		UserID:        req.UserID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Country:       req.Country,
		DeviceID:      req.DeviceID,
		IPAddress:     req.IPAddress,
		Timestamp:     time.Now(),
	}

	eventData, _ := json.Marshal(event)

	err = RDB.XAdd(Ctx, &redis.XAddArgs{
		Stream: "fraud-events-stream",
		Values: map[string]interface{}{
			"data": eventData,
		},
	}).Err()

	if err != nil {
		log.Println("Failed to publish event:", err)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"transaction_id": transactionID,
		"timestamp":      time.Now(),
	})
}