package main

import (
	"context"
	"fmt"

	"github.com/defryheryanto/whatsapp-assistant/internal/assistant"
)

func main() {
	_, err := setupSQLiteConnection("whatsapp_assistant.db")
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	client, err := setupWhatsmeowClient("whatsmeow.db")
	if err != nil {
		panic(fmt.Sprintf("failed to setup whatsmeow client: %v", err))
	}

	whatsappAssistant := assistant.NewWhatsAppAssistant(client)
	whatsappAssistant.Start(context.Background())
}
