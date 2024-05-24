package assistant

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type BirthdayListAction struct {
	*WhatsAppAssistant
}

func (a *BirthdayListAction) Execute(ctx context.Context, evt *events.Message) error {
	chatJid := evt.Info.Chat.ToNonAD().String()

	birthdays, err := a.repository.GetBirthdaysByChatJid(ctx, chatJid)
	if err != nil {
		return err
	}

	output := "Birthday List\n=====================\n"

	for _, bday := range birthdays {
		output += fmt.Sprintf("%s on %s\n", bday.Name, bday.String())
	}

	_, err = a.client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
		Conversation: proto.String(output),
	})
	if err != nil {
		return err
	}

	return nil
}
