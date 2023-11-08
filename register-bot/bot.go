package register_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
)

type Bot struct {
	Token                   string
	bot                     *tgbotapi.BotAPI
	commandHandlers         map[string]func() string
	registerBasicAuth       map[string]int
	registerRegisterCommand map[string]int
	buttons                 map[string]func() string
	privateCommands         map[string]struct{}
	allowedUsernames        map[string]struct{}
	mu                      sync.RWMutex
	onRegisterUser          func(username string) bool // Функция, которая будет вызываться при регистрации пользователя
}

func NewBot(token string, onRegisterUser func(username string) bool) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Token:                   token,
		bot:                     bot,
		commandHandlers:         make(map[string]func() string),
		registerBasicAuth:       make(map[string]int),
		registerRegisterCommand: make(map[string]int),
		privateCommands:         make(map[string]struct{}),
		allowedUsernames:        make(map[string]struct{}),
		mu:                      sync.RWMutex{},
		onRegisterUser:          onRegisterUser, // Устанавливаем функцию обратного вызова
		buttons:                 make(map[string]func() string),
	}, nil
}

func (b *Bot) Start() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		sendKeyboard(b, update.Message.Chat.ID)

		if update.Message != nil {
			if update.Message.IsCommand() {
				command := update.Message.Command()
				handler, ok := b.commandHandlers[command]

				// Флаг, чтобы определить, было ли выполнено действие после if b.registerCommandHandlers[command] == 1
				actionPerformed := true

				if b.registerBasicAuth[command] == 1 {
					// Получаем юзернейм
					username := update.Message.From.UserName
					// Вызываем функцию обратного вызова при регистрации пользователя
					if b.onRegisterUser != nil {
						// если функция обратного вызова вернула false, то handler не отработает
						if !b.onRegisterUser(username) {
							actionPerformed = false
						}
					}
				}

				if b.registerRegisterCommand[command] == 1 {
					// Получаем юзернейм
					username := update.Message.From.UserName
					// Вызываем функцию обратного вызова при регистрации пользователя
					if b.onRegisterUser != nil {
						b.onRegisterUser(username)
					}
				}

				if ok && actionPerformed {
					if b.isPrivateCommand(command) && !b.isAllowedUser(update.Message.From.UserName) {
						// Команда приватна и пользователь не в списке разрешенных
						continue
					}

					messageText := handler()
					if messageText != "" {
						response := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
						_, err := b.bot.Send(response)
						if err != nil {
							log.Println("Ошибка при отправке сообщения:", err)
						}
					}
				}
			} else {
				buttonName := update.Message.Text

				// Флаг, чтобы определить, было ли выполнено действие после if b.registerCommandHandlers[command] == 1
				actionPerformed := true
				if b.registerBasicAuth[buttonName] == 1 {
					// Получаем юзернейм
					username := update.Message.From.UserName
					// Вызываем функцию обратного вызова при регистрации пользователя
					if b.onRegisterUser != nil {
						// если функция обратного вызова вернула false, то handler не отработает
						if !b.onRegisterUser(username) {
							actionPerformed = false
						}
					}
				}

				if b.registerRegisterCommand[buttonName] == 1 {
					// Получаем юзернейм
					username := update.Message.From.UserName
					// Вызываем функцию обратного вызова при регистрации пользователя
					if b.onRegisterUser != nil {
						b.onRegisterUser(username)
					}
				}

				handler, buttonExists := b.buttons[buttonName]

				if buttonExists && actionPerformed {
					//handler, buttonExists := b.buttons[buttonName]
					if buttonExists {
						messageText := handler()
						if messageText != "" {
							response := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
							_, err := b.bot.Send(response)
							if err != nil {
								log.Println("Ошибка при отправке сообщения:", err)
							}
						}
					}
				}
			}
		}
	}
}

func (b *Bot) RegisterButton(buttonText string, command string, handler func() string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buttons[buttonText] = handler

	if handler != nil {
		b.commandHandlers[command] = handler
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
