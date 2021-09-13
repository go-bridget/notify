package notify

import (
	"github.com/go-bridget/mig/cli"
)

type Options struct {
	JwtSecret string
}

func NewOptions() *Options {
	return (&Options{}).Init()
}

func (o *Options) Init() *Options {
	return o
}

func (o *Options) Bind() *Options {
	return o.BindWithPrefix("notify")
}

func (o *Options) BindWithPrefix(prefix string) *Options {
	p := func(s string) string {
		if prefix != "" {
			return prefix + "-" + s
		}
		return s
	}
	cli.StringVar(&o.JwtSecret, p("jwt-secret"), "default", "JWT token signature secret")
	return o
}
