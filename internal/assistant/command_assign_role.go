package assistant

import (
	"context"
	"fmt"
	"strings"

	whatsmeow_proto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type AssignRoleAction struct {
	*WhatsAppAssistant
	Command string
}

func (a *AssignRoleAction) Execute(ctx context.Context, evt *events.Message) error {
	roleName := a.extractRoleName(evt)
	if roleName == "" {
		return nil
	}

	existingCommand := a.getCommands()[roleName]
	if existingCommand != nil {
		_, err := a.client.SendMessage(ctx, evt.Info.Chat, &whatsmeow_proto.Message{
			Conversation: proto.String(fmt.Sprintf("Role '%s' unavailable", roleName)),
		})
		if err != nil {
			return err
		}
		return nil
	}
	groupJid := evt.Info.Chat.ToNonAD().String()

	existingRole, err := a.repository.FindRole(ctx, roleName, groupJid)
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
		GroupJid:   groupJid,
		MemberJIDs: evt.Message.GetExtendedTextMessage().GetContextInfo().GetMentionedJid(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *AssignRoleAction) extractRoleName(evt *events.Message) string {
	message := getMessage(evt)

	roleName := ""
	words := strings.Split(message, " ")
	for i := 1; i < len(words); i++ {
		if words[i-1] == a.Command {
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
