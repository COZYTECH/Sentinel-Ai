package main

import (
	"log"

	"encoding/json"

	"github.com/google/uuid"
)

type FraudCase struct {
	ID       string
	Summary  string
	Embedding []float64
}

// Simple save function
func SaveFraudCase(summary string, embedding []float64) error {
	id := uuid.New().String()
	embeddingJSON, _ := json.Marshal(embedding)
	_, err := DB.Exec(`
		INSERT INTO fraud_case_memory (id, case_summary, embedding_vector)
		VALUES (?, ?, ?)`, id, summary, embeddingJSON)
	return err
}

// Fetch top N similar cases (dummy similarity for demo)
func GetSimilarCases(topN int) ([]FraudCase, error) {
	rows, err := DB.Query(`SELECT id, case_summary FROM fraud_case_memory LIMIT ?`, topN)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cases []FraudCase
	for rows.Next() {
		var fc FraudCase
		if err := rows.Scan(&fc.ID, &fc.Summary); err != nil {
			continue
		}
		cases = append(cases, fc)
	}
	return cases, nil
}


func main() {
	InitDB()
	InitRedis()

	log.Println("Investigation Memory Service running...")

	// Example: saving a new case
	SaveFraudCase("High-value transfer to offshore account", []float64{0.12, -0.33, 0.55})

	// Example: retrieving top 5 similar cases
	cases, _ := GetSimilarCases(5)
	for _, c := range cases {
		log.Println("Similar case:", c.ID, c.Summary)
	}
}