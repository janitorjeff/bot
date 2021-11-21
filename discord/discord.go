package discord

import (
	"fmt"
	"log"

	"github.com/janitorjeff/bot/commands"
	"github.com/janitorjeff/bot/core"
	"github.com/janitorjeff/bot/utils"

	dg "github.com/bwmarrin/discordgo"
)

type Discord struct {
	*dg.Session
	*dg.MessageCreate
}

func (d *Discord) Parse() *core.Message {
	displayName := d.Member.Nick
	if displayName == "" {
		displayName = d.Author.Username
	}

	author := &core.Author{
		ID:          d.Author.ID,
		Name:        d.Author.Username,
		DisplayName: displayName,
		Mention:     d.Author.Mention(),
	}

	channel := &core.Channel{
		ID:   d.ChannelID,
		Name: d.ChannelID,
	}

	cmd, args, prefix := utils.GetCommandArgsPrefix(d.Content, "!")

	command := &core.Cmd{
		Name:      cmd,
		Args:      args,
		HasPrefix: prefix,
		Raw:       d.Content,
	}

	msg := &core.Message{
		ID:       d.ID,
		Type:     core.Discord,
		Author:   author,
		Channel:  channel,
		Command:  command,
		Commands: commands.Commands,
		Client:   d,
	}

	return msg
}

func (d *Discord) Write(msg interface{}) (*core.Message, error) {
	switch t := msg.(type) {
	case string:
		text := msg.(string)
		lenLim := 500
		lenCnt := func(s string) int { return len(s) }

		if lenLim > lenCnt(text) {
			_, err := d.ChannelMessageSend(d.ChannelID, text)
			return nil, err
		}

		parts := utils.Split(text, lenCnt, lenLim)
		for _, p := range parts {
			d.ChannelMessageSend(d.ChannelID, p)
		}

		return nil, nil
		
	case dg.MessageEmbed:
		// TODO
		return nil, nil
	default:
		return nil, fmt.Errorf("Can't send twitch message of type %v", t)
	}

}

func messageCreate(s *dg.Session, m *dg.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.Bot {
		return
	}

	if len(m.Content) == 0 {
		return
	}

	d := &Discord{s, m}
	msg := d.Parse()

	_, err := msg.CommandRun()
	if err != nil {
		log.Println(err)
		return
	}
}

func Init(token string) error {
	d, err := dg.New("Bot " + token)
	if err != nil {
		return err
	}

	d.AddHandler(messageCreate)

	d.Identify.Intents = dg.MakeIntent(dg.IntentsAll)

	d.State = dg.NewState()
	d.State.MaxMessageCount = 100

	return d.Open()
}
