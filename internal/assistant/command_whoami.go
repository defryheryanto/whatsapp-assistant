package assistant

import (
	"context"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type WhoAmIAction struct {
	*WhatsAppAssistant
}

func (a *WhoAmIAction) Execute(ctx context.Context, evt *events.Message) error {
	senderJid := evt.Info.Sender.ToNonAD().String()

	_, err := a.client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
		Conversation: proto.String(senderJid),
	})
	if err != nil {
		return err
	}

	return nil
}
