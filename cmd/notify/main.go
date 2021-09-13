package main

import (
	"context"
	"os"

	"github.com/apex/log"
	"github.com/go-bridget/mig/cli"

	_ "github.com/go-bridget/notify/internal/logger"
	server "github.com/go-bridget/notify/server/notify"
)

var (
	BuildVersion string
	BuildTime    string
)

func main() {
	app := cli.NewApp("notify")

	for _, command := range server.Commands() {
		app.AddCommand(command.Name, command.Title, command.New)
	}

	app.AddCommand("version", "Print version", func() *cli.Command {
		return &cli.Command{
			Run: func(_ context.Context, _ []string) error {
				log.Info(app.Name)
				log.Infof("build version %s", BuildVersion)
				log.Infof("build time    %s", BuildTime)
				return nil
			},
		}
	})
	if err := app.Run(); err != nil {
		log.WithError(err).Error("Exiting")
		os.Exit(1)
	}
}
