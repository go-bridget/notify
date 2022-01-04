package internal

import (
	"context"
)

// SetIPToContext sets IP value to ctx
func SetIPToContext(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, ipAddressKey, ip)
}

// GetIPFromContext gets IP value from ctx
func GetIPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value(ipAddressKey).(string); ok {
		return ip
	}
	return ""
}
