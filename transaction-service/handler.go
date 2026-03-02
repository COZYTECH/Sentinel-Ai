package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"transaction_id": transactionID,
		"timestamp":      time.Now(),
	})
}