package main

type RiskResult struct {
	Score             int
	Level             string
	RecommendedAction string
}

func CalculateRisk(amount float64, country string, historicalCountry string, transactionCountLast5Min int, employeeAccess bool) RiskResult {
	score := 0

	if amount > 10000 {
		score += 30
	}
	if country != historicalCountry {
		score += 20
	}
	if transactionCountLast5Min >= 3 {
		score += 25
	}
	if employeeAccess {
		score += 15
	}

	level := "LOW"
	action := "NONE"

	switch {
	case score >= 75:
		level = "HIGH"
		action = "FREEZE_ACCOUNT"
	case score >= 50:
		level = "MEDIUM"
		action = "REVIEW"
	}

	return RiskResult{
		Score:             score,
		Level:             level,
		RecommendedAction: action,
	}
}