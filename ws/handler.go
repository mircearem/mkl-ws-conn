package socket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(c *http.Request) bool { return true },
}

func Handler(s *Server) echo.HandlerFunc {
	return echo.HandlerFunc(func(ctx echo.Context) error {
		conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
		if err != nil {
			return err
		}
		peer := NewPeer(conn, s.ctx)
		s.addPeerCh <- peer

		return nil
	})
}
