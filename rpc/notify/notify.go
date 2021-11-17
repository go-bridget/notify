package notify

import (
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/twitchtv/twirp"

	"github.com/go-bridget/notify/internal"
	"github.com/go-bridget/notify/internal/websocket"
)

func NewNotifyServiceGateway(svc NotifyService, hooks *twirp.ServerHooks) chi.Router {
	fs := http.FileServer(http.Dir("public_html"))
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Handle("/*", fs)
	mux.Mount("/twirp", NewNotifyServiceServer(svc, hooks))
	mux.HandleFunc("/notify/ws", func(w http.ResponseWriter, req *http.Request) {
		// produce original response writer so we can hijack connection
		if mw, ok := w.(middleware.WrapResponseWriter); ok {
			w = mw.Unwrap()
		}

		var (
			ctx = req.Context()

			// auth result
			res *AuthResponse
		)

		start := func() error {
			// fill out request values
			values, err := internal.FillFromRequest(req, nil)
			if err != nil {
				return err
			}

			// authenticate session
			payload := new(AuthRequest)
			payload.Authorization = values.Get("authorization")

			// authenticate
			res, err = svc.Auth(ctx, payload)
			if err != nil {
				return err
			}
			return nil
		}

		err := start()

		// we now have res.UserID int64
		handler := NewHandler(ctx, res.UserID)
		if err != nil {
			log.WithError(err).Info("Authentication error")
			_ = handler.Send(internal.NewError(err))
			handler.Close()
		} else {
			_ = handler.Send(struct {
				Kind   string `json:"kind"`
				UserID string `json:"userID"`
			}{
				Kind:   "connect",
				UserID: fmt.Sprint(res.UserID),
			})
		}

		ws := websocket.New()
		if err := ws.Open(ctx, w, req, handler); err != nil {
			log.WithError(err).Info("Closing websocket channel")
		}
	})
	return mux
}
