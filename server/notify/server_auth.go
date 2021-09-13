package notify

import (
	"context"

	"github.com/go-bridget/notify/internal/auth"
	rpc "github.com/go-bridget/notify/rpc/notify"
)

var emptyAuthResponse = new(rpc.AuthResponse)

// Auth checks if user can be authenticated
func (svc *Server) Auth(ctx context.Context, r *rpc.AuthRequest) (*rpc.AuthResponse, error) {
	authenticator := auth.NewUserAuthenticator(svc.options.JwtSecret)
	userID, err := authenticator.UserID(r.Authorization)
	if err != nil {
		return emptyAuthResponse, err
	}
	return &rpc.AuthResponse{
		UserID: userID,
	}, nil
}
