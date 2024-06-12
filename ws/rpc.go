package socket

import "io"

type RPC struct {
	From    *Peer
	Payload io.Reader
}
