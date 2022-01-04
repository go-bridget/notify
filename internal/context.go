package internal

type contextKey int

const (
	cookiesKey contextKey = 1 + iota
	ipAddressKey
)
