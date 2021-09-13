package internal

import (
	"net"
	"os"
	"strconv"
	"strings"

	"math/rand"
	"net/http"

	"go.elastic.co/apm/module/apmhttp"
)

// WrapAll wraps a http.Handler with all needed handlers for our service
func WrapAll(h http.Handler) http.Handler {
	h = WrapWithIP(h)
	h = WrapWithAPM(h)
	return h
}

// WrapWithAPM wraps a http.Handler to inject the APM agent
func WrapWithAPM(h http.Handler) http.Handler {
	apmHandler := apmhttp.Wrap(h)
	value := os.Getenv("ELASTIC_APM_GLOBAL_SAMPLE_RATE")
	if value == "" {
		value = "1"
	}
	sampleRate, _ := strconv.ParseFloat(value, 64)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sampleRate > rand.Float64() {
			apmHandler.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// WrapWithIP wraps a http.Handler to inject the client IP into the context
func WrapWithIP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get IP address
		ip := func() string {
			headers := []string{
				http.CanonicalHeaderKey("X-Forwarded-For"),
				http.CanonicalHeaderKey("X-Real-IP"),
			}
			for _, header := range headers {
				if addr := r.Header.Get(header); addr != "" {
					if idx := strings.Index(addr, ","); idx > 0 {
						return addr[:idx]
					}
					return addr
				}
			}
			return r.RemoteAddr
		}()

		// strip ipv6 mapped ipv4 prefix
		ip = strings.TrimPrefix(ip, "::ffff:")

		// Set up addr for sanitizing
		addr := ip

		// RemoteAddr is usually [host]:[port], throw away [port]
		// We use net.SplitHostPort so we can deal with ipv6
		if strings.Contains(ip, ":") {
			var err error
			addr, _, err = net.SplitHostPort(ip)
			if err != nil {
				addr = ip
			}
		}

		ctx := r.Context()
		ctx = SetIPToContext(ctx, addr)

		r.RemoteAddr = addr

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
