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
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// select {
	// case <-ctx.Done():
	// 	reply := tgbotapi.NewMessage(message.Chat.ID, "Превышено время ожидания")
	// 	h.api.Send(reply)
	// 	return
	// default:
	reply := tgbotapi.NewMessage(message.Chat.ID, "start")
	h.api.Send(reply)

}

func (h *BotHandler) HandleNewAppointmentCommand(message *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// select {
	// case <-ctx.Done():
	// 	reply := tgbotapi.NewMessage(message.Chat.ID, "Превышено время ожидания")
	// 	h.api.Send(reply)
	// 	return
	// default:
	result := h.addAppointment(ctx, message)
	reply := tgbotapi.NewMessage(message.Chat.ID, result)
	h.api.Send(reply)

}
