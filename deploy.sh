BINARY_NAME="whatsapp-assistant"
LOG_FILE="./logs/main.log"

# migrate database
go run ./db/migrate.go  >> $LOG_FILE 2>&1 &

# kill & remove the binary
pkill -f "whatsapp-assistant"
rm $BINARY_NAME

# build go binary
go build -o $BINARY_NAME ./cmd/app/...

# run the binary
./$BINARY_NAME >> $LOG_FILE 2>&1 &
