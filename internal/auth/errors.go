package auth

import (
	"github.com/pkg/errors"
)

var (
	errInvalidToken  = errors.New("invalid token")
	errInvalidClaims = errors.New("invalid claims")

	errEmptyToken  = errors.New("empty token")
	errEmptySecret = errors.New("empty secret")
)
