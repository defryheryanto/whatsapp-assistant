package assistant

import (
	"context"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type GetSavedTextAction struct {
	*WhatsAppAssistant
	Command string
}

func (a *GetSavedTextAction) Execute(ctx context.Context, evt *events.Message) error {
	title := a.extractTitle(evt)
	if title == "" {
		return nil
	}

	groupJid := evt.Info.Chat.ToNonAD().String()
	savedText, err := a.repository.GetSavedText(ctx, groupJid, title)
	if err != nil {
		return err
	}
	if savedText == nil {
		return nil
	}

	_, err = a.client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
		Conversation: proto.String(savedText.Content),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *GetSavedTextAction) extractTitle(evt *events.Message) string {
	message := getMessage(evt)
	words := strings.FieldsFunc(message, func(r rune) bool {
		return r == ' ' || r == '\n'
	})

	title := ""
	for i, word := range words {
		if word == a.Command {
			title = words[i+1]
			break
		}
	}

	return title
}
