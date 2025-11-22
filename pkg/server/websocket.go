package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

// WebsocketHandler handles WebSocket connections (currently unused)
// nolint:unused // Preserved for future WebSocket implementation
func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	socket, err := websocket.Accept(w, r, nil)

	if err != nil {
		slog.Error("could not open websocket", "error", err)
		http.Error(w, "could not open websocket", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := socket.Close(websocket.StatusGoingAway, "server closing websocket"); err != nil {
			slog.Error("error closing websocket", "err", err)
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
