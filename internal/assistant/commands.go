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
	IsPrivate   bool
}

const (
	COMMAND_PREFIX        = '%'
	COMMAND_COMMANDS      = "commands"
	COMMAND_ASSIGN_ROLE   = "assign"
	COMMAND_CALL_ROLE     = ""
	COMMAND_CALL_EVERYONE = "all"
	COMMAND_SAVE_TEXT     = "save"
	COMMAND_GET_TEXT      = "text"
	COMMAND_SAVE_BIRTHDAY = "birthday"
	COMMAND_LIST_BIRTHDAY = "birthdaylist"
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
			IsPrivate:   false,
		},
		COMMAND_ASSIGN_ROLE: {
			Format:      fmt.Sprintf("%s [role name] [@member1 @member2 @member3 ...]", commandFormat(COMMAND_ASSIGN_ROLE)),
			Description: "Assign role to mentioned members. To add mention to yourself, use @self",
			Action: &AssignRoleAction{
				WhatsAppAssistant: wa,
				Command:           commandFormat(COMMAND_ASSIGN_ROLE),
			},
			IsPrivate: false,
		},
		COMMAND_CALL_ROLE: {
			Format:      fmt.Sprintf("%s[role name]", commandFormat(COMMAND_CALL_ROLE)),
			Description: "Mention members of called role",
			Action: &CallRoleAction{
				WhatsAppAssistant: wa,
				CommandPrefix:     string(COMMAND_PREFIX),
			},
			IsPrivate: false,
		},
		COMMAND_CALL_EVERYONE: {
			Format:      commandFormat(COMMAND_CALL_EVERYONE),
			Description: "Mention all members in group",
			Action: &CallEveryoneAction{
				WhatsAppAssistant: wa,
			},
			IsPrivate: false,
		},
		COMMAND_SAVE_TEXT: {
			Format:      fmt.Sprintf("%s [title] [content]", commandFormat(COMMAND_SAVE_TEXT)),
			Description: "Save the text",
			Action: &SaveTextAction{
				WhatsAppAssistant: wa,
				Command:           commandFormat(COMMAND_SAVE_TEXT),
			},
			IsPrivate: false,
		},
		COMMAND_GET_TEXT: {
			Format:      fmt.Sprintf("%s [title]", commandFormat(COMMAND_GET_TEXT)),
			Description: "Get the saved text by title",
			Action: &GetSavedTextAction{
				WhatsAppAssistant: wa,
				Command:           commandFormat(COMMAND_GET_TEXT),
			},
			IsPrivate: false,
		},
		COMMAND_SAVE_BIRTHDAY: {
			Format:      fmt.Sprintf("%s [name] [birthday (yyyy-mm-dd)]", commandFormat(COMMAND_SAVE_BIRTHDAY)),
			Description: "Save the person's birthday",
			Action: &PremiumUserAuthenticator{
				client:     wa.client,
				repository: wa.repository,
				next: &SaveBirthdayAction{
					WhatsAppAssistant: wa,
					Command:           commandFormat(COMMAND_SAVE_BIRTHDAY),
				},
			},
			IsPrivate: true,
		},
		COMMAND_LIST_BIRTHDAY: {
			Format:      commandFormat(COMMAND_LIST_BIRTHDAY),
			Description: "List saved birthdays",
			Action: &BirthdayListAction{
				WhatsAppAssistant: wa,
			},
			IsPrivate: false,
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
