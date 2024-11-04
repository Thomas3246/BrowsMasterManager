package handle

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *BotHandler) addAppointment(message *tgbotapi.Message) (resultMessage string) {

	err := h.service.BookingService.CreateAppointment(message.From.ID)

	resultMessage = "Запись успешно добавлена"
	if err != nil {
		resultMessage = "Не удалось создать запись"
	}

	return resultMessage
}
