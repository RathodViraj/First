package notificationservice

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WSService struct {
	Clients map[int]*websocket.Conn
	MU      sync.Mutex
}

var Upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
