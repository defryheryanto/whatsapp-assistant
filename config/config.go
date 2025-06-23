package config

var (
	DatabaseConnectionString string
)

func Init() {
	DatabaseConnectionString = getString("DATABASE_CONNECTION_STRING", "")
}
