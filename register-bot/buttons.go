package register_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (b *Bot) RegisterButton(buttonText string, command string, handler ...func() string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buttons[buttonText] = CommandHandlers{
		TextHandlers: handler,
	}

	if handler != nil {
		b.commandHandlers[command] = CommandHandlers{
			TextHandlers: handler,
		}
	}
}

func sendKeyboard(bot *Bot, chatID int64) {
	// Создайте клавиатуру с кнопками на основе зарегистрированных кнопок
	msg := tgbotapi.NewMessage(chatID, "123")

	if len(bot.buttons) > 0 {
		var keyboardRows [][]tgbotapi.KeyboardButton
		for buttonText := range bot.buttons {
			keyboardRows = append(keyboardRows, []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(buttonText)})
		}

		keyboard := tgbotapi.NewReplyKeyboard(keyboardRows...)
		msg.ReplyMarkup = keyboard
	}

	// сообщение с клавиатурой
	_, err := bot.bot.Send(msg)
	if err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}
}
