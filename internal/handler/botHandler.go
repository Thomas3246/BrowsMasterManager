package handler

import (
	"context"
	"time"

	"github.com/Thomas3246/BrowsMasterManager/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	api     *tgbotapi.BotAPI
	service *service.BotService
}

func NewBotHandler(api *tgbotapi.BotAPI, service *service.BotService) *BotHandler {
	return &BotHandler{
		api:     api,
		service: service}
}

func (h *BotHandler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.api.GetUpdatesChan(u)

	for update := range updates {
		go h.HandleMessage(update)
	}
}

func (h *BotHandler) HandleMessage(update tgbotapi.Update) {
	switch update.Message.Command() {

	case "start":
		h.HandleStartCommand(update.Message)

	case "newAppointment":
		h.HandleNewAppointmentCommand(update.Message)
	}

}

func (h *BotHandler) HandleStartCommand(message *tgbotapi.Message) {

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

func (h *BotHandler) HandleNewAppointmentCommand(message *tgbotapi.Message) {
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
