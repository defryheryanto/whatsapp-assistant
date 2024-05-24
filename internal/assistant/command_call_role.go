package assistant

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type CallRoleAction struct {
	*WhatsAppAssistant
	CommandPrefix string
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
	_, err = a.client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: message,
			ContextInfo: &waE2E.ContextInfo{
				MentionedJID: roles.MemberJIDs,
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
	words := strings.FieldsFunc(message, func(r rune) bool {
		return r == '\n' || r == ' '
	})
	for _, word := range words {
		if string(word[0]) == a.CommandPrefix {
			roleName = word[1:]
			break
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
