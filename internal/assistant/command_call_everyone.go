package assistant

import (
	"context"
	"fmt"
	"strings"

	whatsmeow_proto "go.mau.fi/whatsmeow/binary/proto"
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

	msg := &whatsmeow_proto.Message{
		ExtendedTextMessage: &whatsmeow_proto.ExtendedTextMessage{
			Text: proto.String(strings.Join(mentionText, " ")),
			ContextInfo: &whatsmeow_proto.ContextInfo{
				MentionedJid: mentionedJid,
			},
		},
	}

	_, err = a.client.SendMessage(ctx, evt.Info.Chat, msg)
	if err != nil {
		return err
	}

	return nil
}
