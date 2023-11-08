package register_bot

func (b *Bot) RegisterTextCommand(command string, handler func() string) {
	b.commandHandlers[command] = handler
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
