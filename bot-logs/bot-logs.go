package bot_logs

import (
	"log"
	"net/http"
	"net/url"
)

type TelegramLogger struct {
	chatID string
}

func NewTelegramLogger(chatID string) *TelegramLogger {
	return &TelegramLogger{
		chatID: chatID,
	}
}

func (tl *TelegramLogger) Log(logMessage string) {
	apiURL := "https://api.telegram.org/bot/sendMessage"
	values := url.Values{
		"chat_id": {tl.chatID},
		"text":    {logMessage},
	}

	_, err := http.PostForm(apiURL, values)
	if err != nil {
		log.Printf("Ошибка при отправке лога в телеграм: %v", err)
	}
}
