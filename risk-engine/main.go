package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/COZYTECH/Sentinel-Ai/shared/events"
	"github.com/redis/go-redis/v9"
)

func main() {
	InitDB()
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
			log.Println("Error reading stream:", err)
			time.Sleep(time.Second)
			continue
		}

		for _, msg := range xs[0].Messages {
			lastID = msg.ID

			dataJSON := msg.Values["data"].(string)
			var txEvent events.TransactionCreatedEvent
			if err := json.Unmarshal([]byte(dataJSON), &txEvent); err != nil {
				log.Println("Failed to unmarshal:", err)
				continue
			}

			// TODO: Fetch historical info + employee access
			risk := CalculateRisk(txEvent.Amount, txEvent.Country, txEvent.Country, 0, false)

			log.Printf("Transaction %s Risk Score: %d Level: %s Action: %s\n", txEvent.TransactionID, risk.Score, risk.Level, risk.RecommendedAction)

			// TODO: Save to MySQL risk_assessments table
			// TODO: Publish risk.assessed event
		}
	}
}