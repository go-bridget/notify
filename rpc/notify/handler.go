package notify

import (
	"context"
	"fmt"
	"sync"

	"encoding/json"

	"github.com/go-bridget/notify/internal/pubsub"
)

type Handler struct {
	sync.Mutex
	context.Context

	channel *pubsub.PubSub

	userID   string
	outgoing chan []byte
}

type StatePayload struct {
	Kind  string            `json:"kind"`
	State map[string]string `json:"state"`
}

func NewHandler(ctx context.Context, id string) *Handler {
	return &Handler{
		Context:  ctx,
		channel:  pubsub.New(),
		userID:   id,
		outgoing: make(chan []byte, 512),
	}
}

func (h *Handler) Connect() error {
	// send initial state to user when connecting
	onStart := func() error {
		values, err := h.channel.HGetAll(h, fmt.Sprintf("notify:%s:state", h.userID))
		if err != nil {
			return err
		}
		return h.Send(StatePayload{"state", values})
	}
	// send individual notification messages to user when received
	onMessage := func(channel string, msg []byte) error {
		return h.SendRaw(msg)
	}
	return h.channel.Subscribe(h, fmt.Sprintf("notify:%s", h.userID), onStart, onMessage)
}

func (h *Handler) Close() {
	// clean up when exiting
	h.Lock()
	defer h.Unlock()
	if h.outgoing != nil {
		close(h.outgoing)
		h.outgoing = nil
	}
}

func (*Handler) Incoming(msg []byte) error {
	return nil
}

func (h *Handler) Send(msg interface{}) error {
	h.Lock()
	defer h.Unlock()

	// check if still alive
	if h.outgoing == nil {
		return nil
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return h.SendRaw(b)
}

func (h *Handler) SendRaw(msg []byte) error {
	h.outgoing <- msg
	return nil
}

func (h *Handler) Outgoing() (<-chan []byte, error) {
	return h.outgoing, nil
}
