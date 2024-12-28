package handler

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *BotHandler) handleStartCommand(message *tgbotapi.Message) {

	startMsg := `Привет! Это бот для записи на брови к мастеру ___Евгении___ по адресу **г.Черкесск, ул.**

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

	// Создать функцию registerUser в handler, в ней получать обрабатывать полученное сообщение, выделять из него номер телефона,
	// передавать выделенное в функцию registerUser в service.
	// В service, вызывается функция repository, где проверяется, зарегестрирован ли пользователь с данным номером.
	// Если нет, то должно запрашиваться имя, создаваться сущность пользователя и добавляться в базу.
	// Если да, то у пользователя уточняется его имя (Елена, верно?). Выходят кнопки с вариантами "Да" и "Нет, изменить"

	// После регистрации должно выходить сообщение с инструкциями к пользованию ботом и командами. Должна быть команда /help (она и должна вызываться).
	// В сообщении команды должы быть inline кнопки для записи или т.д.

	// Добавить контексты
}

func (h *BotHandler) handleContact(message *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		reply := tgbotapi.NewMessage(message.Chat.ID, "Превышено время ожидания")
		h.api.Send(reply)
		return
	default:
		contact := message.Contact
		if contact.UserID == message.From.ID {
			// Добавить логику принятия номера - вызов функции сверки с базой
			// Вызов функции уточнения имени

			userName := h.checkForUser(ctx, contact)
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

				h.api.Send(reply)
			} else {
				resultMessage := h.registerUser(ctx, contact)
				reply := tgbotapi.NewMessage(message.Chat.ID, resultMessage)
				reply.ParseMode = "markdown"
				reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				h.api.Send(reply)
			}

		} else {
			msg := "Пожалуйста, отправьте именно ваш контакт, используя кнопку \"**Поделиться**\"."
			reply := tgbotapi.NewMessage(message.Chat.ID, msg)
			reply.ParseMode = "markdown"
			h.api.Send(reply)
		}
	}

}

func (h *BotHandler) handleNewAppointmentCommand(message *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// resultChan := make(chan string, 1)
	// go func() {
	// 	result := h.addAppointment(ctx, message)
	// 	resultChan <- result
	// }()

	select {
	case <-ctx.Done():
		reply := tgbotapi.NewMessage(message.Chat.ID, "Превышено время ожидания")
		h.api.Send(reply)
		return
	// case result := <-resultChan:
	default:
		result := h.addAppointment(ctx, message)
		reply := tgbotapi.NewMessage(message.Chat.ID, result)
		h.api.Send(reply)
	}
}

func (h *BotHandler) handleNameChangeCommand(message *tgbotapi.Message) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		reply := tgbotapi.NewMessage(message.Chat.ID, "Превышено время ожидания")
		h.api.Send(reply)
		return

	default:
		args := message.CommandArguments()
		if args == "" {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Для смены имени необходимо после команды /name указать свое имя. Например: \"/name Лена\"")

			h.api.Send(msg)
			return
		}

		resultMessage := h.changeUserName(ctx, message)
		reply := tgbotapi.NewMessage(message.Chat.ID, resultMessage)
		h.api.Send(reply)
		// добавить MW
	}
}
