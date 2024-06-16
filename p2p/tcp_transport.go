package p2p

import (
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddress string
	listner       net.Listener

	// common practise to put
	// mutex above related map
	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listnerAddr string) *TCPTransport {
	return &TCPTransport{listenAddress: listnerAddr}
}
