package assistant

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

type Command struct {
	Format      string
	Description string
	Action      commandAction
}

const (
	COMMAND_PREFIX        = '%'
	COMMAND_COMMANDS      = "commands"
	COMMAND_ASSIGN_ROLE   = "assign"
	COMMAND_CALL_ROLE     = ""
	COMMAND_CALL_EVERYONE = "all"
)

type commandAction interface {
	Execute(ctx context.Context, evt *events.Message) error
}

func (wa *WhatsAppAssistant) getCommands() map[string]*Command {
	return map[string]*Command{
		COMMAND_COMMANDS: {
			Format:      commandFormat(COMMAND_COMMANDS),
			Description: "Get All Commands",
			Action:      &GetCommandsAction{wa},
		},
		COMMAND_ASSIGN_ROLE: {
			Format:      fmt.Sprintf("%s [role name] [@member1 @member2 @member3 ...]", commandFormat(COMMAND_ASSIGN_ROLE)),
			Description: "Assign role to mentioned members",
			Action: &AssignRoleAction{
				WhatsAppAssistant: wa,
				Command:           commandFormat(COMMAND_ASSIGN_ROLE),
			},
		},
		COMMAND_CALL_ROLE: {
			Format:      fmt.Sprintf("%s[role name]", commandFormat(COMMAND_CALL_ROLE)),
			Description: "Mention members of called role",
			Action: &CallRoleAction{
				WhatsAppAssistant: wa,
				CommandPrefix:     string(COMMAND_PREFIX),
			},
		},
		COMMAND_CALL_EVERYONE: {
			Format:      commandFormat(COMMAND_CALL_EVERYONE),
			Description: "Mention all members in group",
			Action: &CallEveryoneAction{
				WhatsAppAssistant: wa,
			},
		},
	}
}

func (wa *WhatsAppAssistant) getCommandAction(command string) commandAction {
	result := wa.getCommands()[command]
	if result == nil {
		result = wa.getCommands()[COMMAND_CALL_ROLE]
	}
	return result.Action
}

func commandFormat(command string) string {
	return fmt.Sprintf("%c%s", COMMAND_PREFIX, command)
}
