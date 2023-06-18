package assistant

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/types/events"
)

type AssignRoleAction struct {
	*WhatsAppAssistant
}

func (a *AssignRoleAction) Execute(ctx context.Context, evt *events.Message) error {
	roleName := a.extractRoleName(evt)
	if roleName == "" {
		return nil
	}

	existingRole, err := a.repository.FindRole(ctx, roleName)
	if err != nil {
		return err
	}
	if existingRole != nil {
		err = a.repository.DeleteRole(ctx, existingRole.Name)
		if err != nil {
			return err
		}
	}

	err = a.repository.InsertRole(ctx, &Role{
		Name:       roleName,
		MemberJIDs: evt.Message.GetExtendedTextMessage().GetContextInfo().GetMentionedJid(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *AssignRoleAction) extractRoleName(evt *events.Message) string {
	command := fmt.Sprintf("%c%s", COMMAND_PREFIX, COMMAND_ASSIGN_ROLE)
	message := getMessage(evt)

	roleName := ""
	words := strings.Split(message, " ")
	for i := 1; i < len(words); i++ {
		if words[i-1] == command {
			roleName = words[i]
		}
	}

	if roleName == "" {
		return ""
	}
	if roleName[0] == '@' {
		return ""
	}

	return roleName
}
