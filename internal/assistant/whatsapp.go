package assistant

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
)

type Roles struct {
	Name       string
	MemberJIDs []string
}

type WhatsAppAssistantRepository interface {
	InsertRole()
}

type WhatsAppAssistant struct {
	client *whatsmeow.Client
}

func NewWhatsAppAssistant(client *whatsmeow.Client) *WhatsAppAssistant {
	return &WhatsAppAssistant{client}
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	wa.client.Disconnect()
}
