package api

import (
	"context"

	"github.com/labstack/echo/v4"
	ws "github.com/mircearem/mklwsconn/ws"
)

type Server struct {
	app      *echo.Echo
	wsServer *ws.Server
}

func NewServer() *Server {
	ctx := context.Background()
	wsServer := ws.NewServer(ctx)

	app := echo.New()
	app.GET("ws", ws.Handler(wsServer))

	return &Server{
		app:      app,
		wsServer: wsServer,
	}
}

func (s *Server) Run() error {
	go s.wsServer.Start()
	return s.app.Start(":3000")
}
