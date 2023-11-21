package register_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
)

type Bot struct {
	Token                   string
	bot                     *tgbotapi.BotAPI
	commandHandlers         map[string]CommandHandlers
	registerBasicAuth       map[string]int
	registerRegisterCommand map[string]int
	buttons                 map[string]CommandHandlers
	privateCommands         map[string]struct{}
	allowedUsernames        map[string]struct{}
	mu                      sync.RWMutex
	onRegisterUser          func(username string) bool // Функция, которая будет вызываться при регистрации пользователя
	gWorkers                int                        // Количество горутин в пуле для обработки сообщений
}

type CommandHandlers struct {
	TextHandlers  []func() string
	ImageHandlers []func() tgbotapi.Chattable
}

func NewBot(gWorkers int, token string, onRegisterUser func(username string) bool) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Token:                   token,
		bot:                     bot,
		commandHandlers:         make(map[string]CommandHandlers),
		registerBasicAuth:       make(map[string]int),
		registerRegisterCommand: make(map[string]int),
		privateCommands:         make(map[string]struct{}),
		allowedUsernames:        make(map[string]struct{}),
		mu:                      sync.RWMutex{},
		onRegisterUser:          onRegisterUser, // Устанавливаем функцию обратного вызова
		buttons:                 make(map[string]CommandHandlers),
		gWorkers:                gWorkers,
	}, nil
}

func (b *Bot) Start() {
	const workers = 10 // Количество горутин в пуле
	var wg sync.WaitGroup
	jobs := make(chan tgbotapi.Update)

	// Запуск workers горутин для обработки обновлений
	for i := 0; i < b.gWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for update := range jobs {
				b.processUpdate(update)
			}
		}()
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		jobs <- update
	}

	close(jobs)
	wg.Wait()
}

func (b *Bot) processUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	sendKeyboard(b, update.Message.Chat.ID)

	if update.Message.IsCommand() {
		command := update.Message.Command()
		handlers, ok := b.commandHandlers[command]

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
				return
			}

			for _, textHandler := range handlers.TextHandlers {
				messageText := textHandler()
				if messageText != "" {
					response := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
					_, err := b.bot.Send(response)
					if err != nil {
						log.Println("Ошибка при отправке сообщения:", err)
					}
				}
			}
		}
	} else {
		buttonName := update.Message.Text

		// Флаг, чтобы определить, было ли выполнено действие после if b.registerCommandHandlers[command] == 1
		actionPerformedB := true
		if b.registerBasicAuth[buttonName] == 1 {
			// Получаем юзернейм
			username := update.Message.From.UserName
			// Вызываем функцию обратного вызова при регистрации пользователя
			if b.onRegisterUser != nil {
				// если функция обратного вызова вернула false, то handler не отработает
				if !b.onRegisterUser(username) {
					actionPerformedB = false
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

		handlersB, buttonExists := b.buttons[buttonName]

		if buttonExists && actionPerformedB {
			for _, textHandler := range handlersB.TextHandlers {
				messageText := textHandler()
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
