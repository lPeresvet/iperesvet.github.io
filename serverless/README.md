# Serverless Backend for Research Data Collection

## Overview

Yandex Cloud HTTP Function that receives and stores game session results.

## Deploy to Yandex Cloud

### 1. Создать функцию

```bash
cd serverless
zip -r function.zip main.go go.mod

# Создать функцию
yc serverless function create \
  --name research-data-collector \
  --runtime go136 \
  --memory 256MB \
  --execution-timeout 5s \
  --region yandexcloud.net

# Загрузить код
yc serverless functions uploading \
  --function-id <function-id> \
  --zip-archive function.zip
```

### 2. Создать API Gateway

```bash
yc api gateway create \
  --name research-api

# Прикрепить функцию к gateway
yc serverless functions attaching \
  --function-id <function-id> \
  --gateway-id <gateway-id>
```

### 3. Получить URL

```bash
yc api gateway get --name research-api
```

Будет вида: `https://<gateway-id>.api.yandex.net`

### 4. Обновить фронтенд

Замените `API_URL` в `wait/results.html` и `no-wait/results.html`:

```javascript
const API_URL = 'https://<gateway-id>.api.yandex.net/api/save-results';
```

## API Endpoint

```
POST /api/save-results
Content-Type: application/json

Response: 200 OK
{ "status": "ok" }
```

## Data Structure

```json
{
  "uid": "550e8400-e29b-41d4-a716-446655440000",
  "condition": "wait",
  "game1_score": 120,
  "game1_correct": 4,
  "game1_total_rounds": 5,
  "game1_accuracy": 80,
  "game1_maxstreak": 3,
  "game1_time_ms": 25000,
  "game2_score": 250,
  "game2_correct": 8,
  "game2_total_rounds": 10,
  "game2_accuracy": 80,
  "game2_maxstreak": 5,
  "game2_time_ms": 60000,
  "game3_score": 180,
  "game3_correct": 0,
  "game3_total_rounds": 5,
  "game3_accuracy": 0,
  "game3_time_ms": 50000,
  "jump_total_costs": 50,
  "jump_play_count": 3,
  "total_score": 450,
  "timestamp": "2026-10-18T12:00:00Z"
}
```

## Fields

| Field | Type | Description |
|-------|------|-------------|
| uid | string | Unique session identifier (auto-generated UUID) |
| condition | string | Experimental condition: "wait" or "no-wait" |
| game1_score | int | Score from Game 1 (Categories) |
| game1_correct | int | Number of correct answers in Game 1 |
| game1_total_rounds | int | Total rounds in Game 1 (5) |
| game1_accuracy | int | Accuracy percentage (0-100) |
| game1_maxstreak | int | Maximum correct streak in Game 1 |
| game1_time_ms | int | Total time spent on Game 1 (milliseconds) |
| game2_score | int | Score from Game 2 (Odd One Out) |
| game2_correct | int | Number of correct answers in Game 2 |
| game2_total_rounds | int | Total rounds in Game 2 (10) |
| game2_accuracy | int | Accuracy percentage (0-100) |
| game2_maxstreak | int | Maximum correct streak in Game 2 |
| game2_time_ms | int | Total time spent on Game 2 (milliseconds) |
| game3_score | int | Score from Game 3 (Shelf Game) |
| game3_correct | int | Number of correct answers in Game 3 (always 0) |
| game3_total_rounds | int | Total rounds in Game 3 (5) |
| game3_accuracy | int | Accuracy percentage (always 0) |
| game3_time_ms | int | Total time spent on Game 3 (milliseconds) |
| jump_total_costs | int | Total points spent on jump mini-game |
| jump_play_count | int | Number of times jump mini-game was played |
| total_score | int | Final total score (games - jump costs) |
| timestamp | string | ISO 8601 timestamp when data was saved |

## Notes

- Game 3 (Shelf Game) doesn't track accuracy, so `game3_correct`, `game3_accuracy` are always 0
- Yandex Cloud Functions have a 30-second execution limit and 256MB memory limit by default
