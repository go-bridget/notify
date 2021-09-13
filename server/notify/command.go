package notify

import (
	"github.com/go-bridget/mig/cli"
)

func Commands() []*cli.CommandInfo {
	return []*cli.CommandInfo{
		commandStart(),
	}
}
