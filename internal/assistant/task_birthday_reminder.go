package assistant

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mau.fi/whatsmeow"
	whatsmeow_proto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type BirthdayReminder struct {
	client     *whatsmeow.Client
	repository WhatsAppAssistantRepository
}

func NewBirthdayReminder(client *whatsmeow.Client, repository WhatsAppAssistantRepository) *BirthdayReminder {
	return &BirthdayReminder{client, repository}
}

func (r *BirthdayReminder) RunInBackground(ctx context.Context, done chan bool) {
	log.Println("Running birthday reminder task...")
	delay := (24 * time.Hour)

	go func() {
		now := time.Now()
		nextRun := now.Truncate(delay).Add(delay)
		duration := nextRun.Sub(now)

		time.Sleep(duration)
		ticker := time.NewTicker(delay)
		go r.Run(ctx)

		go func() {
			for {
				select {
				case <-done:
					log.Println("reminder stopped")
					ticker.Stop()
					return
				case <-ticker.C:
					go r.Run(ctx)
				}
			}
		}()
	}()
}

func (r *BirthdayReminder) Run(ctx context.Context) error {
	log.Println("checking for birthday reminders...")
	now := time.Now()
	birthdays, err := r.repository.GetBirthdays(ctx, int(now.Month()), now.Year())
	if err != nil {
		log.Printf("error getting birthday: %s\n", err.Error())
		return err
	}

	for _, birthday := range birthdays {
		chatJid, err := types.ParseJID(birthday.TargetChatJid)
		if err != nil {
			log.Printf("error parsing jid %s: %s\n", birthday.TargetChatJid, err.Error())
			continue
		}
		_, err = r.client.SendMessage(ctx, chatJid, &whatsmeow_proto.Message{
			Conversation: proto.String(r.getBirthdayTemplate(birthday)),
		})
		if err != nil {
			log.Printf("error sending message: %s\n", err.Error())
		}
	}

	return nil
}

func (r *BirthdayReminder) getBirthdayTemplate(birthday *Birthday) string {
	return fmt.Sprintf("Today is %s birthday!", birthday.Name)
}
