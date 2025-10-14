package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

// WebsocketHandler handles WebSocket connections (currently unused)
// TODO: Remove if not needed or implement WebSocket functionality
// nolint:unused // Preserved for future WebSocket implementation
func (s *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	socket, err := websocket.Accept(w, r, nil)

	if err != nil {
		log.Printf("could not open websocket: %v", err)
		http.Error(w, "could not open websocket", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := socket.Close(websocket.StatusGoingAway, "server closing websocket"); err != nil {
			log.Printf("error closing websocket: %v", err)
		}
	}()

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		err := socket.Write(socketCtx, websocket.MessageText, []byte(payload))
		if err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
}
