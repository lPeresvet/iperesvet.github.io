package main

import (
	"context"
	"encoding/json"
	"os"
)

type SessionResult struct {
	UID              string `json:"uid"`
	Condition        string `json:"condition"`
	Game1Score       int    `json:"game1_score"`
	Game1Correct     int    `json:"game1_correct"`
	Game1TotalRounds int    `json:"game1_total_rounds"`
	Game1Accuracy    int    `json:"game1_accuracy"`
	Game1MaxStreak   int    `json:"game1_maxstreak"`
	Game1TimeMs      int64  `json:"game1_time_ms"`
	Game2Score       int    `json:"game2_score"`
	Game2Correct     int    `json:"game2_correct"`
	Game2TotalRounds int    `json:"game2_total_rounds"`
	Game2Accuracy    int    `json:"game2_accuracy"`
	Game2MaxStreak   int    `json:"game2_maxstreak"`
	Game2TimeMs      int64  `json:"game2_time_ms"`
	Game3Score       int    `json:"game3_score"`
	Game3Correct     int    `json:"game3_correct"`
	Game3TotalRounds int    `json:"game3_total_rounds"`
	Game3Accuracy    int    `json:"game3_accuracy"`
	Game3TimeMs      int64  `json:"game3_time_ms"`
	JumpTotalCosts   int    `json:"jump_total_costs"`
	JumpPlayCount    int    `json:"jump_play_count"`
	TotalScore       int    `json:"total_score"`
	Timestamp        string `json:"timestamp"`
}

type HTTPEvent struct {
	HTTPMethod string            `json:"httpMethod"`
	Headers    map[string]string `json:"headers"`
	Path       string            `json:"path"`
	Query      string            `json:"query"`
	Body       string            `json:"body"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

func HandleHTTP(ctx context.Context, event HTTPEvent) (Response, error) {
	corsHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}

	if event.HTTPMethod == "OPTIONS" {
		return Response{
			StatusCode: 200,
			Body:       "",
			Headers:    corsHeaders,
		}, nil
	}

	var result SessionResult
	if err := json.Unmarshal([]byte(event.Body), &result); err != nil {
		return Response{StatusCode: 400, Body: `{"error":"invalid json"}`, Headers: corsHeaders}, nil
	}

	if result.UID == "" {
		return Response{StatusCode: 400, Body: `{"error":"missing uid"}`, Headers: corsHeaders}, nil
	}
	if result.Condition != "wait" && result.Condition != "no-wait" {
		return Response{StatusCode: 400, Body: `{"error":"invalid condition"}`, Headers: corsHeaders}, nil
	}

	data, _ := json.Marshal(result)
	filePath := "/function/storage/data/results.jsonl"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return Response{StatusCode: 500, Body: `{"error":"file error"}`, Headers: corsHeaders}, nil
	}
	defer file.Close()

	if _, err := file.Write(append(data, '\n')); err != nil {
		return Response{StatusCode: 500, Body: `{"error":"write error"}`, Headers: corsHeaders}, nil
	}

	return Response{
		StatusCode: 200,
		Body:       `{"status":"ok","uid":"` + result.UID + `"}`,
		Headers:    corsHeaders,
	}, nil
}
