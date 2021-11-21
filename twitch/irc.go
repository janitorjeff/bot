package twitch

import (
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/janitorjeff/bot/commands"
	"github.com/janitorjeff/bot/core"
	"github.com/janitorjeff/bot/utils"

	twitchIRC "github.com/gempir/go-twitch-irc/v2"
)

type TwitchIRC struct {
	*twitchIRC.Client
	*twitchIRC.PrivateMessage
}

var twitchIrcClient *twitchIRC.Client

func (tirc *TwitchIRC) Parse() *core.Message {
	author := &core.Author{
		ID:          tirc.User.ID,
		Name:        tirc.User.Name,
		DisplayName: tirc.User.DisplayName,
		Mention:     fmt.Sprintf("@%s", tirc.User.DisplayName),
	}

	channel := &core.Channel{
		ID:   tirc.RoomID,
		Name: tirc.Channel,
	}

	cmd, args, prefix := utils.GetCommandArgsPrefix(tirc.Message, "!")

	command := &core.Cmd{
		Name:      cmd,
		Args:      args,
		HasPrefix: prefix,
		Raw:       tirc.Message,
	}

	msg := &core.Message{
		ID:       tirc.ID,
		Type:     core.Twitch,
		Author:   author,
		Channel:  channel,
		Command:  command,
		Commands: commands.Commands,
		Client:   tirc,
	}

	return msg
}

func (tirc *TwitchIRC) Write(msg interface{}) (*core.Message, error) {
	var text string
	switch t := msg.(type) {
	case string:
		text = msg.(string)
	default:
		return nil, fmt.Errorf("Can't send twitch message of type %v", t)
	}

	// This is how twitch's server seems to count the length, even though the
	// chat client on twitch's website doesn't follow this
	lenLim := 500
	lenCnt := utf8.RuneCountInString

	if lenLim > lenCnt(text) {
		tirc.Say(tirc.Channel, text)
		return nil, nil
	}

	parts := utils.Split(text, lenCnt, lenLim)
	for _, p := range parts {
		tirc.Say(tirc.Channel, p)
	}

	return nil, nil
}

// func (tirc *TwitchIRC) Delete() error {
// 	_, err := tirc.Write(fmt.Sprintf("/delete %s", tirc.ID))
// 	return err
// }

// func (tirc *TwitchIRC) Edit(msg interface{}) (*core.Message, error) {
// 	return nil, fmt.Errorf("editing not supported for twitch irc")
// }

func onPrivateMessage(m twitchIRC.PrivateMessage) {
	tirc := &TwitchIRC{twitchIrcClient, &m}
	msg := tirc.Parse()

	_, err := msg.CommandRun()
	if err != nil {
		log.Println(err)
		return
	}
}

func IRCInit(nick string, oauth string, channels []string) *twitchIRC.Client {
	twitchIrcClient = twitchIRC.NewClient(nick, oauth)

	twitchIrcClient.OnPrivateMessage(onPrivateMessage)

	twitchIrcClient.Join(channels...)

	go twitchIrcClient.Connect()

	return twitchIrcClient
}
