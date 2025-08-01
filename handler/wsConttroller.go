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

	uidVal, exists := ctx.Get("userID")
	var userID int
	if !exists {
		userID = 25 // or any test user ID
		log.Println("userID not found in context, using default:", userID)
	} else {
		userID, _ = uidVal.(int)
	}

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
