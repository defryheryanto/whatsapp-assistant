package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

func setupWhatsmeowClient(sqlitePath string) (*whatsmeow.Client, error) {
	container, err := sqlstore.New("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", sqlitePath), nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing sqlite: %v", err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("error getting device: %v", err)
	}
	client := whatsmeow.NewClient(deviceStore, nil)
	return client, nil
}
