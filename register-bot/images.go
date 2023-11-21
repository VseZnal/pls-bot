package register_bot

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) SendImageFromBytes(chatID int64, imageData []byte, caption string) error {
	// Создаем буфер
	var buf bytes.Buffer

	// Записываем данные
	buf.Write(imageData)

	// Создаем сообщение для отправки
	photoConfig := tgbotapi.FileBytes{
		Name:  "image.jpg", // Имя файла
		Bytes: buf.Bytes(), // Данные изображения
	}

	// Отправляем сообщение
	msg := tgbotapi.NewPhotoUpload(chatID, photoConfig)
	msg.Caption = caption

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) SendImage(chatID int64, imageFilePath string, caption string) error {
	photoConfig := tgbotapi.NewPhotoUpload(chatID, imageFilePath)
	photoConfig.Caption = caption

	_, err := b.bot.Send(photoConfig)
	return err
}
