package handle

import (
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
		if update.Message != nil {
			h.HandleMessage(update.Message)
		}
	}
}

func (h *BotHandler) HandleMessage(message *tgbotapi.Message) {
	switch message.Command() {

	case "start":

	case "newHandle":
		msg := h.addAppointment(message)
		reply := tgbotapi.NewMessage(message.Chat.ID, msg)
		h.api.Send(reply)
	}

}
