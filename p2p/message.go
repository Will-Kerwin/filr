package p2p

const (
	IncomingStream  = 0x2
	IncomingMessage = 0x1
)

// RPC holds any arbitrary data that is being sent over
// each transport between two nodes in the network
type RPC struct {
	Payload []byte
	From    string
	Stream  bool
}
