package assistant

import (
	"context"
	"log"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func (wa *WhatsAppAssistant) handleCommands(ctx context.Context) whatsmeow.EventHandler {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			message := v.Message.GetConversation()
			log.Printf("Received a message: %s\n", message)
			commands := extractCommands(message)

			for _, command := range commands {
				action := getCommandAction(command)
				if action != nil {
					action(ctx, wa.client, v)
				}
			}
		}
	}
}

func extractCommands(message string) []string {
	words := strings.Split(message, " ")
	commands := []string{}
	for _, word := range words {
		log.Println(word)
		if word[0] == commandPrefix {
			commands = append(commands, word[1:])
		}
	}

	return commands
}
