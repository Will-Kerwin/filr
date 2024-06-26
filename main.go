package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/will-kerwin/filr/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_filr_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {

	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()
	go func() {
		log.Fatal(s2.Start())
	}()

	time.Sleep(time.Second * 2)

	for i := 0; i < 10; i++ {
		data := bytes.NewReader([]byte("my big data file here lol"))
		s2.Store(fmt.Sprintf("myprivatedata_%d", i), data)
		time.Sleep(time.Millisecond * 5)
	}

	// r, err := s2.Get("myprivatedata")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// b, err := io.ReadAll(r)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(b))

	select {}
}
