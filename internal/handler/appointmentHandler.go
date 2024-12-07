package handler

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *BotHandler) addAppointment(ctx context.Context, message *tgbotapi.Message) (resultMessage string) {

	err := h.service.AppointmentService.CreateAppointment(ctx, message.From.ID)

	resultMessage = "Запись успешно добавлена"
	if err != nil {
		resultMessage = "Не удалось создать запись"
		log.Print(err)
	}

	return resultMessage
}
