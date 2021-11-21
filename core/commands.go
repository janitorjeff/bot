package core

import (
	"fmt"
)

type Command struct {
	Name    string
	Aliases []string
	Help    string
	Run     func(*Message) (interface{}, error)
}

type Commands []Command

func (cmds Commands) MatchCommand(cmdName string) (Command, error) {
	if cmdName == "" {
		return Command{}, fmt.Errorf("no command provided")
	}

	for _, c := range cmds {
		if c.Name == cmdName {
			return c, nil
		}

		for _, a := range c.Aliases {
			if a == cmdName {
				return c, nil
			}
		}
	}

	return Command{}, fmt.Errorf("command '%s' not found", cmdName)
}
