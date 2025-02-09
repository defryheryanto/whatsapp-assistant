package assistant

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
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
		_, err := a.client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
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

	mentionedJIDs := evt.Message.GetExtendedTextMessage().GetContextInfo().GetMentionedJID()
	if a.isSelfMention(evt) {
		mentionedJIDs = append(mentionedJIDs, evt.Info.Sender.ToNonAD().String())
	}

	err = a.repository.InsertRole(ctx, &Role{
		Name:       roleName,
		GroupJid:   groupJid,
		MemberJIDs: mentionedJIDs,
	})
	if err != nil {
		return err
	}

	_, err = a.client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
		Conversation: proto.String(fmt.Sprintf("Assigning '%s' success", roleName)),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *AssignRoleAction) extractRoleName(evt *events.Message) string {
	message := getMessage(evt)

	roleName := ""
	words := strings.FieldsFunc(message, func(r rune) bool {
		return r == '\n' || r == ' '
	})
	for i := 1; i < len(words); i++ {
		if words[i-1] == a.Command {
			roleName = words[i]
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

func (a *AssignRoleAction) isSelfMention(evt *events.Message) bool {
	message := getMessage(evt)

	words := strings.FieldsFunc(message, func(r rune) bool {
		return r == '\n' || r == ' '
	})

	for _, word := range words {
		if strings.Contains(word, "@self") {
			return true
		}
	}

	return false
}
