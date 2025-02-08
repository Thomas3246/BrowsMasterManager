package handler

import (
	"context"
	"log"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
)

func (h *BotHandler) addAppointment(ctx context.Context, id int64, appointment *entites.Appointment) (resultMessage string) {

	err := h.service.AppointmentService.CreateAppointment(ctx, id, appointment)

	resultMessage = "Запись успешно добавлена\n\nВы можете просмотреть свои активные записи, нажав на кнопку \"Мои записи\"\n\nИли отменить свою запись, нажав на кнопку \"Отменить запись\""
	if err != nil {
		resultMessage = "Не удалось создать запись\n\nПожалуйста, попробуйте позже"
		log.Print(err)
	}

	return resultMessage
}
