package assistant

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

type GetCommandsAction struct {
	*WhatsAppAssistant
}

func (a *GetCommandsAction) Execute(ctx context.Context, evt *events.Message) error {
	message := "List of available commands (use '%' for command prefix)\n\n"

	for _, command := range a.getCommands() {
		message += fmt.Sprintf("%s: %s\n", command.Format, command.Description)
	}

	_, err := a.client.SendMessage(ctx, evt.Info.Chat, &proto.Message{
		Conversation: &message,
	})
	if err != nil {
		return err
	}

	return nil
}
