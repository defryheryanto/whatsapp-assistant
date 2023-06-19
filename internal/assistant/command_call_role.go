package assistant

import (
	"context"
	"fmt"
	"strings"

	whatsmeow_proto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type CallRoleAction struct {
	*WhatsAppAssistant
	Command string
}

func (a *CallRoleAction) Execute(ctx context.Context, evt *events.Message) error {
	roleName := a.extractRoleName(evt)
	groupJid := evt.Info.Chat.ToNonAD().String()

	roles, err := a.repository.FindRole(ctx, roleName, groupJid)
	if err != nil {
		return err
	}

	mentionText := make([]string, len(roles.MemberJIDs))
	for i, jid := range roles.MemberJIDs {
		phoneNumber := strings.Split(jid, "@")[0]
		mentionText[i] = fmt.Sprintf("@%s", phoneNumber)
	}

	message := proto.String(strings.Join(mentionText, " "))
	_, err = a.client.SendMessage(ctx, evt.Info.Chat, &whatsmeow_proto.Message{
		ExtendedTextMessage: &whatsmeow_proto.ExtendedTextMessage{
			Text: message,
			ContextInfo: &whatsmeow_proto.ContextInfo{
				MentionedJid: roles.MemberJIDs,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *CallRoleAction) extractRoleName(evt *events.Message) string {
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
