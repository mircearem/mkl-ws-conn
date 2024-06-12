package socket

import (
	"context"
	"fmt"
	"log"
	"net"
)

type Server struct {
	peers     map[net.Addr]*Peer
	ctx       context.Context
	errch     chan error
	msgch     chan RPC
	addPeerCh chan *Peer
	delPeerCh chan *Peer
}

func NewServer(ctx context.Context) *Server {
	return &Server{
		peers:     make(map[net.Addr]*Peer),
		ctx:       ctx,
		errch:     make(chan error, 100),
		msgch:     make(chan RPC, 100),
		addPeerCh: make(chan *Peer),
		delPeerCh: make(chan *Peer),
	}
}

func (s *Server) Start() error {
	for {
		select {
		case peer := <-s.addPeerCh:
			s.peers[peer.NetAddr()] = peer
			go s.handlePeer(peer)
			log.Printf("peer: (%+v) connected", peer.NetAddr())
		case peer := <-s.delPeerCh:
			delete(s.peers, peer.conn.RemoteAddr())
		case rpc := <-s.msgch:
			log.Println(rpc)
		case err := <-s.errch:
			log.Println(err.Error())
		}
	}
}

func (s *Server) handlePeer(p *Peer) {
	defer func() {
		p.Close()
		s.delPeerCh <- p
	}()

	errch := make(chan error)

	go func() {
		for rpc := range p.Consume() {
			s.msgch <- rpc
		}
		s.errch <- fmt.Errorf("peer: (%+v) disconnected", p.NetAddr())
	}()

	for {
		select {
		case <-p.Done():
			return
		case <-errch:
			return
		}
	}
}
