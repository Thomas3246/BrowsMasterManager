package handler

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *BotHandler) handleStartCommand(update *tgbotapi.Update) {
	message := update.Message

	startMsg := `Привет! Это бот для записи на брови к мастеру ___ИмяМастера___ по адресу **г.Город, ул.**

		Для записи необходимо подтверждение номера. Для подтверждения нажмите на кнопку *"Поделиться"*.`

	reply := tgbotapi.NewMessage(message.Chat.ID, startMsg)
	reply.ParseMode = "markdown"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("Поделиться"),
		),
	)
	reply.ReplyMarkup = keyboard
	h.api.Send(reply)

	// После регистрации должно выходить сообщение с инструкциями к пользованию ботом и командами. Должна быть команда /help (она и должна вызываться).

}

func (h *BotHandler) handleContact(ctx context.Context, update *tgbotapi.Update) {
	message := update.Message

	resultChan := make(chan tgbotapi.MessageConfig, 1)
	go func() {
		contact := message.Contact
		if contact.UserID == message.From.ID {
			userName, isRegistred, err := h.CheckForUser(ctx, update)
			if err != nil {
				reply := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка, попробуйте позже")
				resultChan <- reply
				return
			}

			if isRegistred {
				if userName != "" {
					msg := "Вы уже зарегистрированы, вас зовут " + userName + ", верно?"

					reply := tgbotapi.NewMessage(message.Chat.ID, msg)
					reply.ParseMode = "markdown"

					inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Да", "callbackConfirmName"),
							tgbotapi.NewInlineKeyboardButtonData("Нет, изменить", "callbackChangeName"),
						),
					)
					reply.ReplyMarkup = inlineKeyboard

					resultChan <- reply
					return
				}

				reply := tgbotapi.NewMessage(message.Chat.ID, "Вы зарегистрированы, но у Вас не указано имя. \nУкажите его через \"/name ___имя___\"")
				reply.ParseMode = "markdown"
				reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				resultChan <- reply
				return

			}
			resultMessage := h.registerUser(ctx, contact)
			reply := tgbotapi.NewMessage(message.Chat.ID, resultMessage)
			reply.ParseMode = "markdown"
			reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			resultChan <- reply

		} else {
			msg := "Пожалуйста, отправьте именно ваш контакт, используя кнопку \"**Поделиться**\"."
			reply := tgbotapi.NewMessage(message.Chat.ID, msg)
			reply.ParseMode = "markdown"
			resultChan <- reply
		}
	}()

	select {
	case <-ctx.Done():
		reply := tgbotapi.NewMessage(message.Chat.ID, "Не удалось обработать\nПревышено время ожидания")
		h.api.Send(reply)
		return

	case result := <-resultChan:
		h.api.Send(result)
	}
}

func (h *BotHandler) handleNewAppointmentCommand(ctx context.Context, update *tgbotapi.Update) {
	message := update.Message

	resultChan := make(chan string, 1)
	go func() {
		result := h.addAppointment(ctx, message)
		resultChan <- result
	}()

	select {
	case <-ctx.Done():
		reply := tgbotapi.NewMessage(message.Chat.ID, "Не удалось создать запись.\nПревышено время ожидания")
		h.api.Send(reply)
		return
	case result := <-resultChan:
		reply := tgbotapi.NewMessage(message.Chat.ID, result)
		h.api.Send(reply)
	}
}

func (h *BotHandler) handleNameChangeCommand(ctx context.Context, update *tgbotapi.Update) {
	message := update.Message

	resultChan := make(chan string, 1)
	go func() {
		result := h.changeUserName(ctx, message)
		resultChan <- result
	}()

	select {
	case <-ctx.Done():
		reply := tgbotapi.NewMessage(message.Chat.ID, "Не удалось изменить имя.\nПревышено время ожидания")
		h.api.Send(reply)
		return

	case result := <-resultChan:
		args := message.CommandArguments()
		if args == "" {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Для смены имени необходимо после команды /name указать свое имя. Например: \"/name Лена\"")

			h.api.Send(msg)
			return
		}

		reply := tgbotapi.NewMessage(message.Chat.ID, result)
		h.api.Send(reply)
	}
}
