package main

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/defryheryanto/whatsapp-assistant/internal/assistant"
	assistant_repository "github.com/defryheryanto/whatsapp-assistant/internal/assistant/repository/gorm"
)

func main() {
	gormDB, err := setupSQLiteConnection(fmt.Sprintf("%s/whatsapp_assistant.db", getAppRootDirectory()))
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	client, err := setupWhatsmeowClient(fmt.Sprintf("%s/whatsmeow.db", getAppRootDirectory()))
	if err != nil {
		panic(fmt.Sprintf("failed to setup whatsmeow client: %v", err))
	}

	whatsappAssistantRepository := assistant_repository.NewWhatsAppAssistantRepository(gormDB)
	whatsappAssistant := assistant.NewWhatsAppAssistant(client, whatsappAssistantRepository)
	whatsappAssistant.Start(context.Background())
}

func getAppRootDirectory() string {
	projectName := regexp.MustCompile(`^(.*whatsapp-assistant)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	return string(rootPath)
}
