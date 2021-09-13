package auth

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type (
	UserAuthenticator struct {
		secret string
	}

	UserClaims struct {
		UserID string `json:"user_id"`
		jwt.StandardClaims
	}
)

func NewUserAuthenticator(secret string) *UserAuthenticator {
	return &UserAuthenticator{
		secret: secret,
	}
}

// UserID retrieves the `user_id` claim from the JWT token
func (u *UserAuthenticator) UserID(token string) (string, error) {
	claims, err := u.UserClaims(token)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// Claims returns the complete JWT claims object
func (u *UserAuthenticator) UserClaims(tokenString string) (*UserClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if tokenString == "" {
		return nil, errEmptyToken
	}
	if u.secret == "" {
		return nil, errEmptySecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errInvalidToken
	}
	if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	}
	return nil, errInvalidClaims
}

// Validate just checks if the JWT claims match an userID
func (u *UserAuthenticator) Validate(token string, userID string) (bool, error) {
	uid, err := u.UserID(token)
	if err != nil {
		return false, err
	}
	return uid == userID, nil
}

// IsUser is a simpler version of Validate, throwing away the error
func (u *UserAuthenticator) IsUser(token string, userID string) bool {
	isUser, _ := u.Validate(token, userID)
	return isUser
}
