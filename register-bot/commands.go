package register_bot

func (b *Bot) RegisterTextCommand(command string, textHandlers ...func() string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.commandHandlers[command] = CommandHandlers{
		TextHandlers: textHandlers,
	}
}

// Регистрация команды для обработки изображения по пути файла
func (b *Bot) RegisterImagePathCommand(command string, imageHandlers ...func() string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.commandHandlers[command] = CommandHandlers{
		ImagePath: imageHandlers,
	}
}

// Регистрация команды для обработки изображения в виде байтов
func (b *Bot) RegisterImageBytesCommand(command string, imageHandlers ...func() []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.commandHandlers[command] = CommandHandlers{
		ImageBytes: imageHandlers,
	}
}

func (b *Bot) BasicAuth(command string) {
	b.registerBasicAuth[command] = 1
}

func (b *Bot) RegisterRegisterCommand(command string) {
	b.registerRegisterCommand[command] = 1
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

func (b *Bot) RegisterUserInputCommand(command string, callback UserInputCallback) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.userInputCallbacks[command] = callback

	b.commandHandlers[command] = CommandHandlers{}
}
