package core

import (
	"fmt"
)

type Author struct {
	ID          string
	Name        string
	DisplayName string
	Mention     string
}

type Channel struct {
	ID   string
	Name string
}

type Cmd struct {
	Name      string
	Args      []string
	HasPrefix bool
	Raw       string
}

type Message struct {
	ID       string
	Type     int
	Command  *Cmd
	Commands Commands
	Author   *Author
	Channel  *Channel
	Client   Platform
}

// Writes a message to the current channel
// returned *Message could be nil depending on the platform
func (m *Message) Write(msg interface{}) (*Message, error) {
	return m.Client.Write(msg)
}

// // Edits a message
// // returned *Message could be nil depending on the platform
// func (m *Message) Edit(msg interface{}) (*Message, error) {
// 	return m.Client.Edit(msg)
// }

// // Deletes a message
// func (m *Message) Delete() error {
// 	return m.Client.Delete()
// }

// Finds, executes and sends a command
func (m *Message) CommandRun() (*Message, error) {
	if !m.Command.HasPrefix {
		return nil, fmt.Errorf("incorrect prefix")
	}

	cmd, err := m.Commands.MatchCommand(m.Command.Name)
	if err != nil {
		return nil, err
	}

	resp, err := cmd.Run(m)
	if err != nil {
		return nil, err
	}

	return m.Write(resp)
}
