package assistant

import (
	"context"
	"strings"

	"go.mau.fi/whatsmeow/types/events"
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

	err := a.repository.SaveText(ctx, &SavedText{
		GroupJid: evt.Info.Chat.ToNonAD().String(),
		Title:    title,
		Content:  content,
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
