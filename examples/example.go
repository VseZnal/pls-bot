package main

import (
	register_bot "github.com/VseZnal/pls-bot/register-bot"
	"io/ioutil"
	"log"
)

func main() {
	// Регистрация ботов
	// Для примера t.me/pls_bot1_bot
	// вторым аргументом можно передать nil, если пользователя не надо обрабатывать
	// 5 - Количество горутин в пуле для обработки сообщений
	bot1, err := register_bot.NewBot(5, "6935692579:AAGZY_RlQceD72lX678YO2FqkLOSig52oLc", func(username string) bool {
		// если будет возвращаться false, то ручка под BasicAuth не отработает, а ручка под RegisterRegisterCommand отработает всегда
		success := GetUser(username)
		return success
	})
	if err != nil {
		log.Fatal(err)
	}

	// Для примера t.me/pls_bot2_bot
	// вторым аргументом можно передать nil, если пользователя не надо обрабатывать
	// 5 - Количество горутин в пуле для обработки сообщений
	bot2, err := register_bot.NewBot(5, "6701968897:AAGLsTyMDBHV_gf5sFxE1XOTXPSBP8kY0Ow", func(username string) bool {
		success := GetUser(username)
		// если будет возвращаться false, то ручка под BasicAuth не отработает, а ручка под RegisterRegisterCommand отработает всегда
		return success
	})
	if err != nil {
		log.Fatal(err)
	}

	////////////////////////////////////////////////////////////////////////////
	// Ручки ///////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////
	// Регистрация хендлеров
	bot1.RegisterTextCommand("text1", handleTextCommand1, handleTextCommand12, handleTextCommand13)
	bot1.RegisterTextCommand("test", handleTextCommand1)
	bot2.RegisterTextCommand("text2", handleTextCommand2)

	// Регистрация команды start будет отправлять сообщение пользователю при первом заходе, если
	// зарегистрированы кнопки, то команда обязательна, так как отправляет и зарегистрированные кнопки пользователю.
	// Зарегистрировать команду возможно только используя RegisterTextCommand
	bot2.RegisterTextCommand("start", handleStartCommand1)

	// под BasicAuth ручка не отработает, если GetUser вернет false
	bot1.BasicAuth("text1")

	// под RegisterRegisterCommand ручка отработает всегда
	bot1.RegisterRegisterCommand("reg")
	bot2.RegisterRegisterCommand("text2")

	// Установка приватности для хендлера
	bot1.SetPrivateCommand("text1")

	// Установка пользователя с правами на приватные методы
	bot1.AllowUser("ZnalZnalZnal")

	// изображения
	bot1.RegisterImageBytesCommand("imageByte", handleImageByte)
	bot1.RegisterImagePathCommand("imagePath", handleImagePath)

	bot1.BasicAuth("imageByte")

	// Регистрация команды, которая будет вызывать функцию обратного вызова
	bot1.RegisterUserInputCommand("get_user_input", handleUserInput)

	////////////////////////////////////////////////////////////////////////////
	// Кнопки //////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////
	bot1.RegisterButton("Кнопка 1", "button1", handleButton1, handleButton12)
	bot2.RegisterButton("Кнопка 2", "button2", handleButton2)

	// под RegisterRegisterCommand ручка отработает всегда
	bot1.RegisterRegisterCommand("Кнопка 1")
	// под BasicAuth ручка не отработает, если GetUser вернет false
	bot2.BasicAuth("Кнопка 2")

	// изображения
	bot1.RegisterButtonImagePathCommand("Кнопка text image", "buttonImagePath", handleButtonImagePath)
	bot1.RegisterButtonImageBytesCommand("Кнопка byte image", "buttonImageByte", handleButtonImageByte)

	// вывод информации по боту
	bot1.PrintRegisteredCommands()
	bot2.PrintRegisteredCommands()

	// Старт бота 1 и бота 2
	go bot1.Start()
	go bot2.Start()

	select {}
}

// Создание функции, которая будет обрабатывать ввод пользователя
func handleUserInput(username, userInput string) string {
	// Здесь можно обрабатывать userInput по вашему усмотрению\
	log.Println(userInput)
	return username + " " + userInput
}

func handleStartCommand1() string {
	return "Это стартовый ответ на команду для бота 2."
}

func handleTextCommand1() string {
	return "Это текстовый ответ на команду для бота 1."
}

func handleTextCommand12() string {
	return "Это текстовый ответ на команду для бота 12."
}

func handleTextCommand13() string {
	return "Это текстовый ответ на команду для бота 13."
}

func handleTextCommand2() string {
	return "Это текстовый ответ на команду для бота 2."
}

func GetUser(username string) bool {
	// обработка юзернейма после register
	log.Println(username)

	// Добавь логику и верни true, если успешно
	return true
}

func handleButton1() string {
	return "Кнопка 1 была нажата."
}

func handleButton12() string {
	return "Кнопка 12 была нажата."
}

func handleButton2() string {
	return "Кнопка 2 была нажата."
}

func handleImageByte() []byte {
	imageData, err := ioutil.ReadFile("./examples/t4k6licnFdc.jpg")
	if err != nil {
		return nil
	}

	return imageData
}

func handleImagePath() string {
	return "./examples/t4k6licnFdc.jpg"
}

func handleButtonImageByte() []byte {
	imageData, err := ioutil.ReadFile("./examples/t4k6licnFdc.jpg")
	if err != nil {
		return nil
	}

	return imageData
}

func handleButtonImagePath() string {
	return "./examples/t4k6licnFdc.jpg"
}
