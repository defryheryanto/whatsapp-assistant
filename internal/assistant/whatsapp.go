package assistant

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

type WhatsAppAssistantRepository interface {
	FindRole(ctx context.Context, name, groupJid string) (*Role, error)
	DeleteRole(ctx context.Context, name string) error
	InsertRole(ctx context.Context, data *Role) error
	SaveText(ctx context.Context, data *SavedText) error
	GetSavedText(ctx context.Context, groupJid, title string) (*SavedText, error)
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
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("WhatsApp Client disconnected")
	wa.client.Disconnect()
}
