package register_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
)

type Bot struct {
	Token            string
	bot              *tgbotapi.BotAPI
	commandHandlers  map[string]func() string
	privateCommands  map[string]struct{}
	allowedUsernames map[string]struct{}
	mu               sync.RWMutex
	onRegisterUser   func(username string) // Функция, которая будет вызываться при регистрации пользователя
}

func NewBot(token string, onRegisterUser func(username string)) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Token:            token,
		bot:              bot,
		commandHandlers:  make(map[string]func() string),
		privateCommands:  make(map[string]struct{}),
		allowedUsernames: make(map[string]struct{}),
		mu:               sync.RWMutex{},
		onRegisterUser:   onRegisterUser, // Устанавливаем функцию обратного вызова
	}, nil
}

func (b *Bot) RegisterTextCommand(command string, handler func() string) {
	b.commandHandlers[command] = handler
}

func (b *Bot) SetPrivateCommand(command string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.privateCommands[command] = struct{}{}
}

func (b *Bot) AllowUser(username string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.allowedUsernames[username] = struct{}{}
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

		if update.Message.IsCommand() {
			command := update.Message.Command()
			handler, ok := b.commandHandlers[command]
			if ok {
				if b.isPrivateCommand(command) && !b.isAllowedUser(update.Message.From.UserName) {
					// Команда приватна и пользователь не в списке разрешенных
					continue
				}

				if command == "register" {
					// Получаем юзернейм
					username := update.Message.From.UserName

					// Вызываем функцию обратного вызова при регистрации пользователя
					if b.onRegisterUser != nil {
						b.onRegisterUser(username)
					}
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
		}
	}
}

func (b *Bot) isPrivateCommand(command string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, isPrivate := b.privateCommands[command]
	return isPrivate
}

func (b *Bot) isAllowedUser(username string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, allowed := b.allowedUsernames[username]
	return allowed
}
