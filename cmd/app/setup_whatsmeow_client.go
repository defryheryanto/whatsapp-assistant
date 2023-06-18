package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func setupWhatsmeowClient(sqlitePath string) (*whatsmeow.Client, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", sqlitePath), dbLog)
	if err != nil {
		return nil, fmt.Errorf("error initializing sqlite: %v", err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("error getting device: %v", err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	return client, nil
}
