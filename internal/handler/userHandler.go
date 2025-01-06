package handler

import (
	"context"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *BotHandler) registerUser(ctx context.Context, contact *tgbotapi.Contact) (resultMessage string) {

	// В регистрацию поменять имя на логин тг. Чтобы при регистрации не было пустого имени, а пустое имя == пользователь не зареган

	// или сделать сессию регистрации при отправке номера, чтобы сразу указывалось имя

	phone := contact.PhoneNumber
	id := strconv.FormatInt(contact.UserID, 10)

	err := h.service.UserService.RegisterUser(ctx, id, phone)
	resultMessage = "Вы были успешно зарегистрированы. Пожалуйста, укажите ваше имя командой \n\"/name ___имя___\". \nНапример: \"/name Лена\""
	if err != nil {
		resultMessage = "Произошла ошибка регистрации"
		log.Println(err)
	}
	return resultMessage
}

func (h *BotHandler) checkForUser(ctx context.Context, contact *tgbotapi.Contact) (userName string) {

	phone := contact.PhoneNumber

	userName = h.service.UserService.CheckForUser(ctx, phone)

	return userName
}

func (h *BotHandler) changeUserName(ctx context.Context, message *tgbotapi.Message) (resultMessage string) {
	newName := message.CommandArguments()
	id := strconv.FormatInt(message.From.ID, 10)

	err := h.service.UserService.ChangeUserName(ctx, id, newName)
	resultMessage = "Имя было успешно изменено"
	if err != nil {
		resultMessage = "Не удалось изменить имя"
		log.Print(err)
	}

	return resultMessage
}
