package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type (
	Websocket struct {
		// the websocket connection
		*websocket.Conn

		// timeout configs
		timeout     time.Duration
		pingTimeout time.Duration
		pingPeriod  time.Duration
	}
)

func New() *Websocket {
	const (
		timeout     = 15 * time.Second
		pingTimeout = 120 * time.Second
		pingPeriod  = (pingTimeout * 9) / 10
	)
	ws := &Websocket{
		timeout:     timeout,
		pingTimeout: pingTimeout,
		pingPeriod:  pingPeriod,
	}

	return ws
}

// Handles websocket requests from peers
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// Allow connections from any Origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *Websocket) readLoop(handler Handler) (err error) {
	defer func() {
		handler.Close()
	}()

	for {
		_, raw, err := ws.Conn.ReadMessage()
		if err != nil {
			// client cancelled connection
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return nil
			}
			return errors.Wrap(err, "ws.readLoop")
		}

		if err := handler.Incoming(raw); err != nil {
			return errors.Wrap(err, "ws.readLoop Incoming")
		}
	}
}

func (ws *Websocket) write(messageType int, payload []byte) error {
	if err := ws.Conn.SetWriteDeadline(time.Now().Add(ws.timeout)); err != nil {
		return err
	}
	return ws.Conn.WriteMessage(messageType, payload)
}

func (ws *Websocket) writeLoop(handler Handler) error {
	ticker := time.NewTicker(ws.pingPeriod)

	defer func() {
		ticker.Stop()
		handler.Close()
	}()

	outgoing, err := handler.Outgoing()
	if err != nil {
		return errors.Wrap(err, "ws.writeLoop Outgoing")
	}

	for {
		select {
		case msg, ok := <-outgoing:
			// channel closed
			if !ok {
				return nil
			}
			// non-empty message
			if msg != nil {
				if err := ws.write(websocket.TextMessage, msg); err != nil {
					return errors.Wrap(err, "ws.writeLoop")
				}
			}
		case <-handler.Done():
			return nil
		case <-ticker.C:
			if err := ws.write(websocket.PingMessage, nil); err != nil {
				return errors.Wrap(err, "ws.writeLoop ping")
			}
		}
	}
}

func (ws *Websocket) Open(ctx context.Context, w http.ResponseWriter, r *http.Request, handler Handler) error {
	// upgrade http request to ws connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		var wsErr websocket.HandshakeError
		if errors.As(err, &wsErr) {
			return errors.Wrap(err, "ws: need a websocket handshake")
		}
		return errors.Wrap(err, "ws: failed to upgrade connection")
	}

	// set connection ping timeouts
	if err = conn.SetReadDeadline(time.Now().Add(ws.pingTimeout)); err != nil {
		return err
	}
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(ws.pingTimeout))
	})

	ws.Conn = conn

	defer func() {
		handler.Close()
		_ = ws.Conn.Close()
	}()

	finalErr := make(chan error, 1)
	go func() {
		finalErr <- handler.Connect()
	}()
	go func() {
		finalErr <- ws.readLoop(handler)
	}()
	go func() {
		finalErr <- ws.writeLoop(handler)
	}()

	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(err, "http request closed")
		case err := <-finalErr:
			return errors.Wrap(err, "error in read/write")
		}
	}
}
