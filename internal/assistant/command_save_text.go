package assistant

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type SaveTextAction struct {
	*WhatsAppAssistant
	Command string
}

func (a *SaveTextAction) Execute(ctx context.Context, evt *events.Message) error {
	title, content := a.extractTitleAndContent(evt)
	if title == "" || content == "" {
		return nil
	}

	groupJid := evt.Info.Chat.ToNonAD().String()
	existingSavedText, err := a.repository.GetSavedText(ctx, groupJid, title)
	if err != nil {
		return err
	}
	if existingSavedText != nil {
		if err = a.repository.DeleteSavedText(ctx, groupJid, title); err != nil {
			return err
		}
	}

	err = a.repository.SaveText(ctx, &SavedText{
		GroupJid: groupJid,
		Title:    title,
		Content:  content,
	})
	if err != nil {
		return err
	}

	_, err = a.client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
		Conversation: proto.String(fmt.Sprintf("Title '%s' Saved", title)),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *SaveTextAction) extractTitleAndContent(evt *events.Message) (string, string) {
	message := getMessage(evt)
	words := strings.FieldsFunc(message, func(r rune) bool {
		return r == ' ' || r == '\n'
	})

	commandIndex := -1
	for i, word := range words {
		if word == a.Command {
			commandIndex = i
			break
		}
	}

	title := ""
	content := ""
	if commandIndex >= 0 && len(words) > commandIndex+2 {
		title = words[commandIndex+1]
		contentStartIndex := len(strings.Join(words[:commandIndex+2], " ")) + 1
		content = message[contentStartIndex:]
	}

	return title, content
}
