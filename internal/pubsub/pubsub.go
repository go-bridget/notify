package pubsub

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/gomodule/redigo/redis"
)

type PubSub struct {
	addr string

	timeout     time.Duration
	pingTimeout time.Duration
	pingPeriod  time.Duration
}

func New() *PubSub {
	const (
		timeout     = 15 * time.Second
		pingTimeout = 120 * time.Second
		pingPeriod  = (pingTimeout * 9) / 10
	)
	var addr string
	if addr = os.Getenv("PUBSUB_REDIS_ADDR"); addr == "" {
		addr = "redis:6379"
	}
	if !strings.Contains(addr, ":") {
		addr += ":6379"
	}
	return &PubSub{
		addr:        addr,
		timeout:     timeout,
		pingTimeout: pingTimeout,
		pingPeriod:  pingPeriod,
	}
}

func (ps *PubSub) dial() (redis.Conn, error) {
	return redis.Dial(
		"tcp",
		ps.addr,
		redis.DialReadTimeout(ps.pingTimeout+ps.timeout),
		redis.DialWriteTimeout(ps.timeout),
	)
}

func (ps *PubSub) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	// main redis connection
	conn, err := ps.dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return redis.StringMap(conn.Do("HGETALL", key))
}

func (ps *PubSub) Subscribe(ctx context.Context, channel string, onStart func() error, onMessage func(channel string, payload []byte) error) error {
	// main redis connection
	conn, err := ps.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	// pubsub object
	psc := redis.PubSubConn{Conn: conn}
	if err := psc.Subscribe(redis.Args{}.Add(channel)...); err != nil {
		return err
	}

	done := make(chan error, 1)

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if err := onMessage(n.Channel, n.Data); err != nil {
					done <- err
					return
				}
			case redis.Subscription:
				switch n.Count {
				case 1:
					// Notify application when all channels are subscribed.
					if err := onStart(); err != nil {
						done <- err
						return
					}
				case 0:
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()

	defer func() {
		if err := psc.Unsubscribe(); err != nil {
			log.WithError(err).Debug("Unsubscribe from redis pubsub")
		}
	}()

	for {
		select {
		case <-time.After(ps.pingPeriod):
			if err := psc.Ping(""); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		case err := <-done:
			return err
		}
	}
}

func (ps *PubSub) Publish(ctx context.Context, channel, message string) error {
	conn, err := ps.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("PUBLISH", channel, message)
	return err
}
