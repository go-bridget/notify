package logger

import (
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

// Logger is a custom split stdout/stderr cli logger
type Logger struct {
	Stdout log.Handler
	Stderr log.Handler
}

// HandleLog splits log outputs to stdout and stderr
func (l *Logger) HandleLog(e *log.Entry) error {
	// Send error/fatal to stderr
	if e.Level == log.ErrorLevel || e.Level == log.FatalLevel {
		return l.Stderr.HandleLog(e)
	}
	// Everything else (info) to stdout
	return l.Stdout.HandleLog(e)
}

func init() {
	if !strings.HasPrefix(os.Getenv("ENVIRONMENT"), "prod") {
		log.SetLevel(log.DebugLevel)
	}
	log.SetHandler(&Logger{
		Stdout: cli.New(os.Stdout),
		Stderr: cli.New(os.Stderr),
	})
}
