package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/labstack/echo/v4"
)

// websocketHandler handles WebSocket connections (currently unused)
// TODO: Remove if not needed or implement WebSocket functionality
// nolint:unused // Preserved for future WebSocket implementation
func (s *Server) websocketHandler(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()
	socket, err := websocket.Accept(w, r, nil)

	if err != nil {
		log.Printf("could not open websocket: %v", err)
		_, _ = w.Write([]byte("could not open websocket"))
		w.WriteHeader(http.StatusInternalServerError)
		return nil
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
	return nil
}
