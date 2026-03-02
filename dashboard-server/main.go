package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RDB *redis.Client
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	log.Println("Connected to Redis")
}

func wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		// Keep connection alive
		if _, _, err := conn.NextReader(); err != nil {
			delete(clients, conn)
			break
		}
	}
}

func broadcast(message []byte) {
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			delete(clients, conn)
		}
	}
}

func listenRedis() {
	stream := "risk-assessed-stream"
	lastID := "0"
	for {
		xs, err := RDB.XRead(Ctx, &redis.XReadArgs{
			Streams: []string{stream, lastID},
			Block:   0,
			Count:   1,
		}).Result()
		if err != nil {
			continue
		}
		for _, msg := range xs[0].Messages {
			lastID = msg.ID
			dataJSON := msg.Values["data"].(string)
			broadcast([]byte(dataJSON))
		}
	}
}

func main() {
	InitRedis()
	go listenRedis()

	r := gin.Default()
	r.GET("/ws", wsHandler)
	r.StaticFile("/", "./index.html")
	r.Run(":8085")
}