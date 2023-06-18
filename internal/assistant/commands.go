package assistant

import (
	"context"

	"go.mau.fi/whatsmeow/types/events"
)

type Command struct {
	Description string
	Action      commandAction
}

const commandPrefix = '%'

type commandAction interface {
	Execute(ctx context.Context, evt *events.Message) error
}

func (wa *WhatsAppAssistant) getCommands() map[string]*Command {
	return map[string]*Command{
		"commands": {
			Description: "Get All Commands",
			Action:      &GetCommandsAction{wa},
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
