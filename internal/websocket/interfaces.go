package websocket

type Handler interface {
	// Handler start
	Connect() error

	// Handler closure
	Close()

	// An incoming message sink
	Incoming([]byte) error

	// A channel for outgoing messages
	Outgoing() (<-chan []byte, error)

	// Basically a passthru to outer ctx.Done()
	Done() <-chan struct{}
}
