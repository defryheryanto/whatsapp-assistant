package assistant

import (
	"context"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func (wa *WhatsAppAssistant) handleCommands(ctx context.Context) whatsmeow.EventHandler {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			message := getMessage(v)
			commands := extractCommands(message)

			for _, command := range commands {
				action := wa.getCommandAction(command)
				if action != nil {
					action.Execute(ctx, v)
				}
			}
		}
	}
}

func getMessage(evt *events.Message) string {
	if evt.Message.GetConversation() != "" {
		return evt.Message.GetConversation()
	}
	if evt.Message.GetExtendedTextMessage().Text != nil {
		return evt.Message.GetExtendedTextMessage().GetText()
	}
	return ""
}

func extractCommands(message string) []string {
	words := strings.Split(message, " ")
	commands := []string{}
	for _, word := range words {
		if word[0] == COMMAND_PREFIX {
			commands = append(commands, word[1:])
		}
	}

	return commands
}
