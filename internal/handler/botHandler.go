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
		// Для каждого сообщения создается своя горутина для параллелизма
		go h.HandleMessage(&update)
	}
}

func (h *BotHandler) HandleMessage(update *tgbotapi.Update) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Обрабатывается отправка пользователем контакта
	if update.Message != nil {
		if update.Message.Contact != nil {
			h.handleContact(ctx, update)
		} else {
			// Обрабаывается отправка пользователем команд
			switch update.Message.Command() {

			case "start":
				// MiddleWare не требуется, общедоступная команда
				h.handleStartCommand(update)

			case "appointment":
				// Для выполнения записи сначала производится проверка, зарегистрирован ли пользователь,
				// после чего проверяется, ввел ли пользователь свое имя
				handler := h.NameMiddleWare(h.AuthMiddleWare(h.handleNewAppointmentCommand))
				handler(ctx, update)

			case "name":
				// Для изменения имени необходимо подтверждение, что пользователь зарегистрирован
				handler := h.AuthMiddleWare(h.handleNameChangeCommand)
				handler(ctx, update)
			}

		}
	}

	// Обработка нажатий на кнопки
	if update.CallbackQuery != nil {
		callbackQuery := update.CallbackQuery
		switch callbackQuery.Data {
		case "callbackConfirmName":
			h.handleConfirmNameCallback(update)

		case "callbackChangeName":
			h.handleChangeNameCallback(update)
		}

	}
}
