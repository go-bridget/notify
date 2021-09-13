package notify

import (
	"context"
)

func New(ctx context.Context, options *Options) (*Server, error) {
	server := &Server{
		context: ctx,
		options: options,
	}
	return server, nil
}
