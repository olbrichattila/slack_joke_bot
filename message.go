package main

import (
	"os"

	"github.com/slack-go/slack"
)

type MessagerInterface interface {
	Send(string) error
}

type Messager struct {
}

func NewMessager() *Messager {
	return &Messager{}
}

func (m *Messager) Send(message string) error {
	token := os.Getenv("SLACK_BOT_TOKEN")
	channelID := os.Getenv("SLACK_CHANNEL_ID")

	api := slack.New(token)

	msgOptions := slack.MsgOptionText(message, false)

	_, _, err := api.PostMessage(channelID, msgOptions)

	return err

}
