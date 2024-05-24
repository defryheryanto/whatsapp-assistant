package assistant

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var ErrForbidden = fmt.Errorf("sender is not premium user")

// PremiumUserAuthenticator is a decorator for checking whether the user is on our premium users list.
// This authenticator should implement 'commandAction' interface
type PremiumUserAuthenticator struct {
	client     *whatsmeow.Client
	repository WhatsAppAssistantRepository
	next       commandAction
}

func NewPremiumUserAuthenticator(
	client *whatsmeow.Client,
	repository WhatsAppAssistantRepository,
	next commandAction,
) *PremiumUserAuthenticator {
	return &PremiumUserAuthenticator{client, repository, next}
}

func (a *PremiumUserAuthenticator) Execute(ctx context.Context, evt *events.Message) error {
	senderJid := evt.Info.Sender.ToNonAD().String()

	premiumUser, err := a.repository.GetPremiumUser(ctx, senderJid)
	if err != nil {
		return err
	}
	if premiumUser == nil {
		_, err := a.client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
			Conversation: proto.String("You are not a premium user."),
		})
		if err != nil {
			return err
		}
		return nil
	}

	return a.next.Execute(ctx, evt)
}
