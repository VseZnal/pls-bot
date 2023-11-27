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

func (b *Bot) RegisterButtonImagePathCommand(buttonText string, command string, handler ...func() string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buttons[buttonText] = CommandHandlers{
		ImagePath: handler,
	}

	if handler != nil {
		b.commandHandlers[command] = CommandHandlers{
			ImagePath: handler,
		}
	}
}

func (b *Bot) RegisterButtonImageBytesCommand(buttonText string, command string, handler ...func() []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buttons[buttonText] = CommandHandlers{
		ImageBytes: handler,
	}

	if handler != nil {
		b.commandHandlers[command] = CommandHandlers{
			ImageBytes: handler,
		}
	}
}

func sendKeyboard(bot *Bot, chatID int64, startMsg string) {
	// Создаем клавиатуру с кнопками на основе зарегистрированных кнопок
	msg := tgbotapi.NewMessage(chatID, startMsg)

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
