package internal

import (
	"context"
	"net/http"
)

// Cookies reads the cookies from the context
func Cookies(ctx context.Context) ([]*http.Cookie, bool) {
	value, ok := ctx.Value(cookiesKey).([]*http.Cookie)
	return value, ok
}

// WrapWithCookies adds the http.Request cookies into context
func WrapWithCookies(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := setCookies(r.Context(), r.Cookies())
		base.ServeHTTP(w, r.WithContext(ctx))
	})
}

// setCookies creates a context with a cookies value
func setCookies(ctx context.Context, cookies []*http.Cookie) context.Context {
	return context.WithValue(ctx, cookiesKey, cookies)
}
