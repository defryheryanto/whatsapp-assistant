package main

import (
	"context"
	"fmt"

	"github.com/defryheryanto/whatsapp-assistant/config"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"google.golang.org/protobuf/proto"
)

func setupWhatsmeowClient(ctx context.Context, sqlitePath string) (*whatsmeow.Client, error) {
	store.DeviceProps.Os = proto.String(config.AppName)
	container, err := sqlstore.New(ctx, "sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", sqlitePath), nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing sqlite: %v", err)
	}
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting device: %v", err)
	}
	client := whatsmeow.NewClient(deviceStore, nil)
	return client, nil
}
