LOG_FILE="./logs/down_time.log"

# Check if the binary is running
if ! pgrep -f "whatsapp-assistant" > /dev/null; then
    echo "[$(date)] whatsapp-assistant is not running. Starting..." >> $LOG_FILE
    ./deploy.sh
    echo "[$(date)] whatsapp-assistant started." >> $LOG_FILE
else
    # If already running, log the status
    echo "[$(date)] whatsapp-assistant is already running."
fi
