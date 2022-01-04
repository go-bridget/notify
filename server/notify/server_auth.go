package notify

import (
	"context"

	"github.com/go-bridget/notify/internal"
	"github.com/go-bridget/notify/internal/auth"
	rpc "github.com/go-bridget/notify/rpc/notify"
)

var emptyAuthResponse = new(rpc.AuthResponse)

// getSession produces a session ID from a parameter or cookie value
func (svc *Server) getSessionID(ctx context.Context, authorization string, sessionCookie string) string {
	cookies, ok := internal.Cookies(ctx)
	if !ok {
		return authorization
	}

	for _, cookie := range cookies {
		if cookie.Name == sessionCookie {
			return cookie.Value
		}
	}
	return authorization
}

// Auth checks if user can be authenticated
func (svc *Server) Auth(ctx context.Context, r *rpc.AuthRequest) (*rpc.AuthResponse, error) {
	// Produce Session ID from cookie if unset with r.Authorization
	sessionID := svc.getSessionID(ctx, r.Authorization, svc.options.SessionCookie)

	authenticator := auth.NewUserAuthenticator(svc.options.JwtSecret)
	userID, err := authenticator.UserID(sessionID)
	if err != nil {
		return emptyAuthResponse, err
	}
	return &rpc.AuthResponse{
		UserID: userID,
	}, nil
}
