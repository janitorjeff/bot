package commands

import (
	"fmt"

	"github.com/janitorjeff/bot/core"
)

func cmdTestGeneric(m *core.Message) (string, error) {
	return fmt.Sprintf("%s -> Test command", m.Author.Mention), nil
}

func cmdTest(m *core.Message) (interface{}, error) {
	switch m.Type {
	default:
		return cmdTestGeneric(m)
	}
}

var Test = core.Command{
	Name: "test",
	Aliases: []string{
		"alias",
	},
	Help:    "test command",
	Run: cmdTest,
}
