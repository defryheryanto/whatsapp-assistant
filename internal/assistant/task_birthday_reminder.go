package assistant

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
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
	birthdays, err := r.repository.GetBirthdays(ctx, int(now.Day()), int(now.Month()))
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

		message, err := r.getBirthdayMessage(birthday, chatJid)
		if err != nil {
			log.Printf("error getting message: %v\n", err)
		}

		_, err = r.client.SendMessage(ctx, chatJid, message)
		if err != nil {
			log.Printf("error sending message: %s\n", err.Error())
		}
	}

	return nil
}

func (r *BirthdayReminder) getBirthdayMessage(birthday *Birthday, chatJid types.JID) (*waE2E.Message, error) {
	age := time.Now().Year() - int(birthday.BirthYear)
	basicMessage := fmt.Sprintf("%s turning %d today!", birthday.Name, age)

	isGroup, err := r.isGroup(chatJid)
	if err != nil {
		return nil, err
	}
	if !isGroup {
		return &waE2E.Message{
			Conversation: proto.String(basicMessage),
		}, nil
	}

	mentionText, mentionedJid, err := r.mentionAll(chatJid)
	if err != nil {
		return nil, err
	}

	return &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(fmt.Sprintf("%s\n%s", basicMessage, mentionText)),
			ContextInfo: &waE2E.ContextInfo{
				MentionedJID: mentionedJid,
			},
		},
	}, nil
}

func (r *BirthdayReminder) isGroup(chatJid types.JID) (bool, error) {
	joinedGroups, err := r.client.GetJoinedGroups()
	if err != nil {
		return false, err
	}
	for _, gr := range joinedGroups {
		if gr.JID.String() == chatJid.String() {
			return true, nil
		}
	}

	return false, nil
}

func (r *BirthdayReminder) mentionAll(chatJid types.JID) (mentionText string, mentionedJid []string, err error) {
	groupInfo, err := r.client.GetGroupInfo(chatJid)
	if err != nil {
		return "", nil, err
	}

	mentionedNumber := make([]string, len(groupInfo.Participants))
	mentionedJid = make([]string, len(groupInfo.Participants))
	for i, participant := range groupInfo.Participants {
		jid := participant.JID.ToNonAD().String()
		phoneNumber := strings.Split(jid, "@")[0]

		mentionedNumber[i] = fmt.Sprintf("@%s", phoneNumber)
		mentionedJid[i] = jid
	}

	return strings.Join(mentionedNumber, " "), mentionedJid, nil
}
