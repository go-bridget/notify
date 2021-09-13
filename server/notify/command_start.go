package notify

import (
	"context"
	"time"

	"net/http"

	"github.com/apex/log"
	"github.com/go-bridget/mig/cli"
	"github.com/pkg/errors"

	"github.com/go-bridget/notify/internal"
	"github.com/go-bridget/notify/rpc/notify"
)

func commandStart() *cli.CommandInfo {
	var (
		options   *Options
		server    *Server
		err       error
		addr      = internal.Getenv("PORT", ":3000")
		pprofAddr = internal.Getenv("PPROF_PORT", ":6060")
		serverErr = make(chan error, 1)
	)
	return &cli.CommandInfo{
		Name:  "start",
		Title: "Start notify service",
		New: func() *cli.Command {
			return &cli.Command{
				Bind: func(ctx context.Context) {
					go internal.NewMonitor(15)
					options = NewOptions().Bind()
				},
				Init: func(ctx context.Context) error {
					server, err = New(ctx, options)
					if err != nil {
						return err
					}
					return nil
				},
				Run: func(ctx context.Context, commands []string) error {
					if err = server.Start(ctx); err != nil {
						return err
					}
					defer server.Shutdown()

					twirpHandler := notify.NewNotifyServiceGateway(server, internal.NewServerHooks())

					httpServer := &http.Server{
						Addr:           addr,
						Handler:        internal.WrapAll(twirpHandler),
						ReadTimeout:    30 * time.Second,
						WriteTimeout:   30 * time.Second,
						MaxHeaderBytes: http.DefaultMaxHeaderBytes,
					}

					// start a debug server for pprof debugging purposes
					go func() {
						log.Infof("Starting pprof on port %s", pprofAddr)
						_ = http.ListenAndServe(pprofAddr, nil)
					}()

					// start the rpc service
					go func() {
						log.Infof("Starting service on port %s", addr)
						err := httpServer.ListenAndServe()
						if !errors.Is(err, http.ErrServerClosed) {
							serverErr <- err
						}
					}()

					// listen for kill signal or rpc exit
					select {
					case <-ctx.Done():
						err = ctx.Err()
					case err = <-serverErr:
					}

					log.WithError(err).Error("shutting down")
					log.Errorf("#+v", err)

					// attempt graceful shutdown in 5sec, return possible error
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					return httpServer.Shutdown(ctx)
				},
			}
		},
	}
}
