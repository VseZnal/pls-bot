package main

import (
	"log"
	register_bot "pls-bot/register-bot"
)

func main() {
	// Регистрация ботов
	bot1, err := register_bot.NewBot("BOT1_TOKEN", func(username string) {
		GetUser(username)
	})
	if err != nil {
		log.Fatal(err)
	}

	bot2, err := register_bot.NewBot("BOT2_TOKEN", func(username string) {
		GetUser(username)
	})
	if err != nil {
		log.Fatal(err)
	}

	// Регистрация хендлеров
	bot1.RegisterTextCommand("text1", handleTextCommand1)
	bot2.RegisterTextCommand("text2", handleTextCommand2)

	// Установка приватности для хендлера
	bot1.SetPrivateCommand("text1")

	// Установка пользователя с правами на приватные методы
	bot1.AllowUser("ZnalZnalZnal")

	// Старт бота 1 и бота 2
	go bot1.Start()
	go bot2.Start()

	select {}
}

func handleTextCommand1() string {
	return "Это текстовый ответ на команду для бота 1."
}

func handleTextCommand2() string {
	return "Это текстовый ответ на команду для бота 2."
}

func GetUser(username string) {
	// обработка юзернейма после register
}
