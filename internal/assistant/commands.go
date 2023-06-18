package assistant

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

type Command struct {
	Description string
	Action      commandAction
}

const commandPrefix = '%'

type commandAction func(ctx context.Context, evt *events.Message) error

func (wa *WhatsAppAssistant) getCommands() map[string]*Command {
	return map[string]*Command{
		"commands": {
			Description: "Get All Commands",
			Action:      wa.getAvailableCommands,
		},
	}
}

func (wa *WhatsAppAssistant) getCommandAction(command string) commandAction {
	result := wa.getCommands()[command]
	if result == nil {
		return nil
	}
	return result.Action
}

func (wa *WhatsAppAssistant) getAvailableCommands(ctx context.Context, evt *events.Message) error {
	message := "List of available commands (use '%' for command prefix)\n\n"

	for key, command := range wa.getCommands() {
		message += fmt.Sprintf("%c%s: %s\n", commandPrefix, key, command.Description)
	}

	_, err := wa.client.SendMessage(ctx, evt.Info.Chat, &proto.Message{
		Conversation: &message,
	})
	if err != nil {
		return err
	}

	return nil
}
