package main

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type FetchFileRequest struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
}

type discord interface {
	FetchFile(ctx context.Context, req *FetchFileRequest) (string, error)
}

type Discord struct {
	*discordgo.Session

	Token string
}

func NewDiscord(token string) discord {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(fmt.Sprintf("failed to create Discord session: %v", err))
	}

	return &Discord{
		Token:   token,
		Session: session,
	}
}

func (d *Discord) FetchFile(ctx context.Context, req *FetchFileRequest) (string, error) {
	msg, err := d.Session.ChannelMessage(req.ChannelID, req.MessageID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch message: %v", err)
	}

	if len(msg.Attachments) == 0 {
		return "", fmt.Errorf("message has no attachments")
	}

	return msg.Attachments[0].URL, nil
}
