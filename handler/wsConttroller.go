package handler

import (
	"First/model"
	"First/notification"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *notification.Hub, ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	log.Println("WebSocket connection established")

	// For local testing, always use userID 25
	var userID int = 25
	log.Println("For testing, using userID 25 for WebSocket connection")

	client := &notification.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan model.Notification),
		Hub:    hub,
	}

	hub.Register <- client

	log.Printf("Client registered with userID: %d\n", userID)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		client.ReadPump()
		wg.Done()
	}()
	go func() {
		client.WritePump()
		wg.Done()
	}()
	wg.Wait()
}
