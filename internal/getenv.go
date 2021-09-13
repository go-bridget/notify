package internal

import (
	"os"
)

func Getenv(key string, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}
