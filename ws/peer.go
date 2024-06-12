package socket

import (
	"bytes"
	"context"
	"net"

	"github.com/gorilla/websocket"
)

type Peer struct {
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func NewPeer(conn *websocket.Conn, ctx context.Context) *Peer {
	c, cancel := context.WithCancel(ctx)
	return &Peer{
		conn:   conn,
		ctx:    c,
		cancel: cancel,
	}
}

func (p *Peer) NetAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *Peer) Consume() chan RPC {
	rpcch := make(chan RPC)

	go func() {
		for {
			_, msg, err := p.conn.ReadMessage()
			if err != nil {
				break
			}
			rpcch <- RPC{
				From:    p,
				Payload: bytes.NewReader(msg),
			}
		}
		close(rpcch)
		p.cancel()

	}()
	return rpcch
}

func (p *Peer) Write(msg []byte) error {
	return p.conn.WriteJSON(msg)
}

func (p *Peer) Done() <-chan struct{} {
	return p.ctx.Done()
}

func (p *Peer) Close() error {
	return p.conn.Close()
}

func (p *Peer) Cancel() error {
	p.cancel()
	return nil
}
