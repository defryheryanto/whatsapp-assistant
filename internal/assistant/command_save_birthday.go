package assistant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type SaveBirthdayAction struct {
	*WhatsAppAssistant
	Command string
}

func (a *SaveBirthdayAction) Execute(ctx context.Context, evt *events.Message) error {
	name, birthday := a.extractNameAndBirthday(evt)

	birthDate, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return err
	}

	chatJid := evt.Info.Chat.ToNonAD().String()

	existingBirthday, err := a.repository.GetBirthday(ctx, name, chatJid)
	if err != nil {
		return err
	}
	if existingBirthday != nil {
		_, err = a.client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
			Conversation: proto.String(
				fmt.Sprintf(
					"%s's birthday is already set on %s",
					existingBirthday.Name,
					existingBirthday.String(),
				),
			),
		})
		if err != nil {
			return err
		}
		return nil
	}

	newBirthday := &Birthday{
		Name:          name,
		BirthDate:     int16(birthDate.Day()),
		BirthMonth:    int16(birthDate.Month()),
		BirthYear:     int16(birthDate.Year()),
		TargetChatJid: chatJid,
	}
	err = a.repository.InsertBirthday(ctx, newBirthday)
	if err != nil {
		return err
	}

	_, err = a.client.SendMessage(ctx, evt.Info.Chat, &waE2E.Message{
		Conversation: proto.String(fmt.Sprintf(
			"saved %s birthday on %s",
			strings.ToUpper(newBirthday.Name),
			newBirthday.String(),
		)),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *SaveBirthdayAction) extractNameAndBirthday(evt *events.Message) (string, string) {
	message := getMessage(evt)
	words := strings.FieldsFunc(message, func(r rune) bool {
		return r == ' ' || r == '\n'
	})

	commandIndex := -1
	for i, word := range words {
		if word == a.Command {
			commandIndex = i
			break
		}
	}

	name := ""
	birthday := ""
	if commandIndex >= 0 && len(words) > commandIndex+2 {
		name = words[commandIndex+1]
		contentStartIndex := len(strings.Join(words[:commandIndex+2], " ")) + 1
		birthday = message[contentStartIndex:]
	}

	return name, birthday
}
