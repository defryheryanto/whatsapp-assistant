package assistant

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type CallEveryoneAction struct {
	*WhatsAppAssistant
}

func (a *CallEveryoneAction) Execute(ctx context.Context, evt *events.Message) error {
	if !evt.Info.IsGroup {
		return nil
	}

	groupInfo, err := a.client.GetGroupInfo(evt.Info.Chat)
	if err != nil {
		return err
	}

	mentionText := make([]string, len(groupInfo.Participants))
	mentionedJid := make([]string, len(groupInfo.Participants))
	for i, participant := range groupInfo.Participants {
		jid := participant.JID.ToNonAD().String()
		phoneNumber := strings.Split(jid, "@")[0]

		mentionText[i] = fmt.Sprintf("@%s", phoneNumber)
		mentionedJid[i] = jid
	}

	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(strings.Join(mentionText, " ")),
			ContextInfo: &waE2E.ContextInfo{
				MentionedJID: mentionedJid,
			},
		},
	}

	_, err = a.client.SendMessage(ctx, evt.Info.Chat, msg)
	if err != nil {
		return err
	}

	return nil
}
