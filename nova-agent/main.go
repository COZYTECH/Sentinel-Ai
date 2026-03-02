package main

import (
	"encoding/json"
	"log"
	"time"

	"context"

	"github.com/COZYTECH/Sentinel-Ai/shared/events"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	log.Println("Connected to Redis")
}

type NovaOutput struct {
	RiskScore          int      `json:"risk_score"`
	RiskLevel          string   `json:"risk_level"`
	PrimaryRiskFactors []string `json:"primary_risk_factors"`
	BehavioralAnalysis string   `json:"behavioral_analysis"`
	EmployeeRisk       string   `json:"employee_involvement_risk"`
	RecommendedAction  string   `json:"recommended_action"`
}

// Dummy Nova Agent reasoning function
func CallNovaAgent(tx events.TransactionCreatedEvent) NovaOutput {
	// For demo, just simple rules
	score := 50
	level := "MEDIUM"
	action := "REVIEW"

	if tx.Amount > 10000 {
		score = 87
		level = "HIGH"
		action = "FREEZE_ACCOUNT"
	}

	return NovaOutput{
		RiskScore:          score,
		RiskLevel:          level,
		PrimaryRiskFactors: []string{"High transaction amount", "New country"},
		BehavioralAnalysis: "Transaction deviates from user profile",
		EmployeeRisk:       "LOW",
		RecommendedAction:  action,
	}
}

func main() {
	InitRedis()
	stream := "fraud-events-stream"
	lastID := "0"

	for {
		xs, err := RDB.XRead(Ctx, &redis.XReadArgs{
			Streams: []string{stream, lastID},
			Block:   0,
			Count:   1,
		}).Result()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		for _, msg := range xs[0].Messages {
			lastID = msg.ID
			dataJSON := msg.Values["data"].(string)
			var tx events.TransactionCreatedEvent
			if err := json.Unmarshal([]byte(dataJSON), &tx); err != nil {
				continue
			}

			// Call AI reasoning
			output := CallNovaAgent(tx)

			// Publish risk.assessed event
			outJSON, _ := json.Marshal(output)
			RDB.XAdd(Ctx, &redis.XAddArgs{
				Stream: "risk-assessed-stream",
				Values: map[string]interface{}{"data": outJSON},
			})
			log.Println("Nova Agent processed transaction", tx.TransactionID)
		}
	}
}