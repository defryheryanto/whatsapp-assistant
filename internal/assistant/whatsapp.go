package assistant

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mau.fi/whatsmeow"
)

type Role struct {
	Name       string
	GroupJid   string
	MemberJIDs []string
}

type SavedText struct {
	GroupJid string
	Title    string
	Content  string
}

type Birthday struct {
	Name          string
	BirthDate     int16
	BirthMonth    int16
	BirthYear     int16
	TargetChatJid string
}

func (b *Birthday) String() string {
	formattedDate := fmt.Sprintf("%d-%02d-%02d", b.BirthYear, b.BirthMonth, b.BirthDate)

	birthdateString := formattedDate
	birthdayDate, err := time.Parse("2006-01-02", formattedDate)
	if err == nil {
		birthdateString = birthdayDate.Format("January 2, 2006")
	}

	return birthdateString
}

type PremiumUser struct {
	UserJid string
}

type WhatsAppAssistantRepository interface {
	FindRole(ctx context.Context, name, groupJid string) (*Role, error)
	DeleteRole(ctx context.Context, name string) error
	InsertRole(ctx context.Context, data *Role) error
	SaveText(ctx context.Context, data *SavedText) error
	GetSavedText(ctx context.Context, groupJid, title string) (*SavedText, error)
	DeleteSavedText(ctx context.Context, groupJid, title string) error
	InsertBirthday(ctx context.Context, birthday *Birthday) error
	GetBirthdays(ctx context.Context, date, month int) ([]*Birthday, error)
	GetBirthday(ctx context.Context, name, chatJid string) (*Birthday, error)
	GetPremiumUser(ctx context.Context, senderJid string) (*PremiumUser, error)
	GetBirthdaysByChatJid(ctx context.Context, chatJid string) ([]*Birthday, error)
}

type WhatsAppAssistant struct {
	client     *whatsmeow.Client
	repository WhatsAppAssistantRepository
}

func NewWhatsAppAssistant(client *whatsmeow.Client, repository WhatsAppAssistantRepository) *WhatsAppAssistant {
	return &WhatsAppAssistant{client, repository}
}

func (wa *WhatsAppAssistant) Start(ctx context.Context) {
	wa.client.AddEventHandler(wa.handleCommands(ctx))
	if wa.client.Store.ID == nil {
		qrChan, _ := wa.client.GetQRChannel(context.Background())
		err := wa.client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				log.Println("QR code:", evt.Code)
			} else {
				log.Println("Login event:", evt.Event)
			}
		}
	} else {
		err := wa.client.Connect()
		if err != nil {
			panic(err)
		}
	}

	log.Println("WhatsApp Client has connected")
	c := make(chan os.Signal, 1)
	scheduleDoneChan := make(chan bool, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	wa.RunBackgroundTasks(ctx, scheduleDoneChan)
	<-c
	scheduleDoneChan <- true

	log.Println("WhatsApp Client disconnected")
	wa.client.Disconnect()
}
