package notify

import (
	"context"

	"github.com/go-bridget/notify/rpc/notify"
)

// Server implements notify.Notify
type Server struct {
	context context.Context
	options *Options
}

// Start is a start hook after flags parsing
func (*Server) Start(_ context.Context) error {
	return nil
}

// Shutdown is a cleanup hook after SIGTERM
func (*Server) Shutdown() {
}

var _ notify.NotifyService = &Server{}
