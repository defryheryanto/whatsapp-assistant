package assistant

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

type Command struct {
	Description string
	Action      commandAction
}

const commandPrefix = '%'

type commandAction func(ctx context.Context, client *whatsmeow.Client, evt *events.Message) error

func getCommands() map[string]*Command {
	return map[string]*Command{
		"commands": {
			Description: "Get All Commands",
			Action:      getAvailableCommands,
		},
	}
}

func getCommandAction(command string) commandAction {
	return getCommands()[command].Action
}

func getAvailableCommands(ctx context.Context, client *whatsmeow.Client, evt *events.Message) error {
	message := "List of available commands (use '%' for command prefix)\n\n"

	for key, command := range getCommands() {
		message += fmt.Sprintf("%c%s: %s\n", commandPrefix, key, command.Description)
	}

	_, err := client.SendMessage(ctx, evt.Info.Chat, &proto.Message{
		Conversation: &message,
	})
	if err != nil {
		return err
	}

	return nil
}
