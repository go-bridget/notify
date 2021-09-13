package auth

import (
	"os"
	"testing"
	"time"

	"github.com/apex/log"
	"github.com/dgrijalva/jwt-go"
)

func getJwtSecret() string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default"
	}
	return jwtSecret
}

func getJwtUserClaim(userID string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	return claims
}

func getJwt(claims jwt.MapClaims, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return at.SignedString([]byte(secret))
}

func TestAuth(t *testing.T) {
	t.Parallel()

	uid := "2"

	jwtSecret := getJwtSecret()
	jwtClaims := getJwtUserClaim(uid)

	token, err := getJwt(jwtClaims, jwtSecret)
	if err != nil {
		t.Fatal(err)
	}

	aa := NewUserAuthenticator(jwtSecret)
	if !aa.IsUser(token, uid) {
		t.Fatalf("Can't verify claim for user_id=%s, token=%s", uid, token)
	}

	log.Infof("Generated JWT: %s", token)
}
