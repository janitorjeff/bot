package core

const (
	Discord = iota
	Twitch
)

type Platform interface {
	Parse() *Message
	Write(interface{}) (*Message, error)
	// Edit(interface{}) (*Message, error)
	// Delete() error
}
