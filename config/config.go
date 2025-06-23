package config

var (
	AppName                  string
	DatabaseConnectionString string
)

func Init() {
	AppName = getString("APP_NAME", "whatsapp-assistant-dev")
	DatabaseConnectionString = getString("DATABASE_CONNECTION_STRING", "")
}
