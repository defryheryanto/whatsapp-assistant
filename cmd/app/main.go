package main

import "fmt"

func main() {
	_, err := setupSQLiteConnection("whatsapp_assistant.db")
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
}
