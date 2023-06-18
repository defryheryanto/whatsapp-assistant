package assistant

import (
	"context"

	"go.mau.fi/whatsmeow/types/events"
)

type Command struct {
	Description string
	Action      commandAction
}

const (
	COMMAND_PREFIX      = '%'
	COMMAND_COMMANDS    = "commands"
	COMMAND_ASSIGN_ROLE = "assign"
)

type commandAction interface {
	Execute(ctx context.Context, evt *events.Message) error
}

func (wa *WhatsAppAssistant) getCommands() map[string]*Command {
	return map[string]*Command{
		COMMAND_COMMANDS: {
			Description: "Get All Commands",
			Action:      &GetCommandsAction{wa},
		},
		COMMAND_ASSIGN_ROLE: {
			Description: "Assign role to mentioned members",
			Action:      &AssignRoleAction{wa},
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
