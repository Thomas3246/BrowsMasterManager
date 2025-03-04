package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *BotHandler) registerUser(ctx context.Context, contact *tgbotapi.Contact) (resultMessage string) {

	phone := contact.PhoneNumber
	id := strconv.FormatInt(contact.UserID, 10)

	// указать номер мастера или достать его из новой функции
	err := h.service.UserService.RegisterUser(ctx, id, phone)
	resultMessage = "Вы были успешно зарегистрированы. Пожалуйста, укажите ваше имя командой \n\"/name ___имя___\". \nНапример: \"/name Лена\""
	if err != nil {
		resultMessage = "Произошла ошибка регистрации"
		log.Println(err)
	}
	return resultMessage
}

func (h *BotHandler) CheckForUser(ctx context.Context, update *tgbotapi.Update) (userName string, isRegistred bool, err error) {

	userId := strconv.FormatInt(update.SentFrom().ID, 10)

	userName, isRegistred, err = h.service.UserService.CheckForUser(ctx, userId)
	if err != nil {
		log.Println(err)
		return "", isRegistred, err
	}

	return userName, isRegistred, nil
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

func (h *BotHandler) handleMyAppointments(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	appointments, err := h.service.UserService.CheckForAppointments(ctx, update.Message.From.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			noRowsMsg := tgbotapi.NewMessage(update.Message.From.ID, "У вас еще нет записей")
			h.api.Send(noRowsMsg)
			return
		}

		errMsg := tgbotapi.NewMessage(update.Message.From.ID, "Произошла ошибка\nПожалуйста, попробуйте позже")
		h.api.Send(errMsg)
		log.Println("Ошибка определения записей пользователя: ", err)
		return
	}

	if len(appointments) == 0 {
		noRowsMsg := tgbotapi.NewMessage(update.Message.From.ID, "У вас еще нет записей")
		h.api.Send(noRowsMsg)
		return
	}

	msgText := "Ваши записи:"

	for i, appointment := range appointments {
		services, err := h.service.ServiceService.GetServicesInAppointment(ctx, appointments[i].ID)
		if err != nil {
			errMsg := tgbotapi.NewMessage(update.Message.From.ID, "Произошла ошибка\nПожалуйста, попробуйте позже")
			h.api.Send(errMsg)
			log.Println("Ошибка определения сервисов в записи пользователя: ", err)
			return
		}

		now := time.Now()
		now = now.Add(-time.Duration(now.Hour())*time.Hour - time.Duration(now.Minute())*time.Minute)
		hour, _ := strconv.Atoi(appointment.Hour)
		minute, _ := strconv.Atoi(appointment.Minute)
		now = now.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute)
		after := now.Add(time.Duration(appointment.TotalDuration) * time.Minute)
		msgText = fmt.Sprintf(msgText+"\n\n%s\n%s-%s\nУслуги:\n", appointment.DateStr, now.Format("15:04"), after.Format("15:04"))

		for _, service := range services {
			msgText = fmt.Sprintf(msgText+"%s\n", service.Name)
		}

		msgText = fmt.Sprintf(msgText+"Стоимость: %d", appointment.TotalCost)
	}

	msg := tgbotapi.NewMessage(update.Message.From.ID, msgText)
	h.api.Send(msg)
}
