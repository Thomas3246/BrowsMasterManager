package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	rusdate "github.com/Thomas3246/BrowsMasterManager/pkg/rusDate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Кнопка подтверждения имени при уточнении
func (h *BotHandler) handleConfirmNameCallback(update *tgbotapi.Update) {
	callbackQuery := update.CallbackQuery
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	h.api.Request(callback)

	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Отлично, так и оставим")
	attachFunctionalButtons(&msg)
	h.api.Send(msg)
}

// Кнопка смены имени при уточнении
func (h *BotHandler) handleChangeNameCallback(update *tgbotapi.Update) {
	callbackQuery := update.CallbackQuery
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	h.api.Request(callback)

	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Для смены имени напишите команду \"/name ___имя___\". \nНапример: \"/name Лена\"")
	msg.ParseMode = "markdown"
	attachFunctionalButtons(&msg)
	h.api.Send(msg)
}

func (h *BotHandler) handleDateChooseCallback(callbackQuery *tgbotapi.CallbackQuery, dayNumber int, userAppointments *entites.Appointment) {
	date := time.Now()
	date = date.Add(-time.Duration(date.Hour())*time.Hour - time.Duration(date.Minute())*time.Minute - time.Duration(date.Second())*time.Second)

	editText := "Давайте выберем дату\n\n" + rusdate.FormatDayMonth(date.Add(time.Hour*24*time.Duration(dayNumber)))

	editKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			// tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(date.Day()) + date., ""),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date), "date_0"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24)), "date_1"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*2)), "date_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*3)), "date_3"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*4)), "date_4"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*5)), "date_5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(functionalButtons.confirm, "confirmDate"),
		),
	)

	editMsg := tgbotapi.NewEditMessageTextAndMarkup(
		callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID,
		editText,
		editKeyboard,
	)

	userAppointments.Date = date.Add(time.Hour * 24 * time.Duration(dayNumber))

	h.api.Send(editMsg)

}

func (h *BotHandler) handleDateConfirmCallback(callbackQuery *tgbotapi.CallbackQuery, userAppointment *entites.Appointment) {
	editText := "**" + rusdate.FormatDayMonth(userAppointment.Date) + "**\n\nНа какое время?\n\n"

	var rows [4][4]tgbotapi.InlineKeyboardButton

	startTime := userAppointment.Date.Add(12 * time.Hour)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			btnText := startTime.Format("15:04")
			callBackText := "chooseTime:" + btnText

			// Сделать после || Проверку на наличие записи в это время у БД

			if time.Now().After(startTime.Add(-30 * time.Minute)) {
				btnText = "❌" + btnText
				callBackText = "inactiveTime"
			} else {
				btnText = "☑️" + btnText
			}

			rows[i][j] = tgbotapi.InlineKeyboardButton{
				Text:         btnText,
				CallbackData: &callBackText,
			}

			startTime = startTime.Add(30 * time.Minute)
		}
	}

	editKeyboard := tgbotapi.NewInlineKeyboardMarkup(rows[0][:], rows[1][:], rows[2][:], rows[3][:],
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(functionalButtons.back, "backToDate")))

	editMsg := tgbotapi.NewEditMessageTextAndMarkup(
		callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID,
		editText,
		editKeyboard,
	)

	editMsg.ParseMode = "markdown"
	h.api.Send(editMsg)
}

func (h *BotHandler) handleBackToDate(callbackQuery *tgbotapi.CallbackQuery) {
	msgText := "Давайте выберем дату\n\n"

	date := time.Now()
	dateKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			// tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(date.Day()) + date., ""),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date), "date_0"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24)), "date_1"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*2)), "date_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*3)), "date_3"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*4)), "date_4"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*5)), "date_5"),
		),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(
		callbackQuery.From.ID,
		callbackQuery.Message.MessageID,
		msgText,
		dateKeyboard,
	)
	h.api.Send(msg)
}

func (h *BotHandler) handleTimeChooseCallback(callbackQuery *tgbotapi.CallbackQuery, hour string, minute string, userAppointments *entites.Appointment) {
	timeStr := hour + ":" + minute
	// userAppointments.DateTime = userAppointments.DateTime.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute) // put into confirm
	editText := rusdate.FormatDayMonth(userAppointments.Date) + "\n\n" + timeStr + "\n\nНа какое время?\n\n"

	var rows [4][4]tgbotapi.InlineKeyboardButton

	startTime := userAppointments.Date.Add(12 * time.Hour)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			btnText := startTime.Format("15:04")
			callBackText := "chooseTime:" + btnText

			// Сделать после || Проверку на наличие записи в это время у БД

			if time.Now().After(startTime.Add(-30 * time.Minute)) {
				btnText = "❌" + btnText
				callBackText = "inactiveTime"
			} else if btnText == hour+":"+minute {
				btnText = "✅" + btnText
			} else {
				btnText = "☑️" + btnText
			}

			rows[i][j] = tgbotapi.InlineKeyboardButton{
				Text:         btnText,
				CallbackData: &callBackText,
			}

			startTime = startTime.Add(30 * time.Minute)
		}
	}

	editKeyboard := tgbotapi.NewInlineKeyboardMarkup(rows[0][:], rows[1][:], rows[2][:], rows[3][:],
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(functionalButtons.confirm, "confirmTime")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(functionalButtons.back, "backToDate")))

	editMsg := tgbotapi.NewEditMessageTextAndMarkup(
		callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID,
		editText,
		editKeyboard,
	)

	editMsg.ParseMode = "markdown"
	h.api.Send(editMsg)

	userAppointments.Hour = hour
	userAppointments.Minute = minute
}

func (h *BotHandler) handleTimeConfirmCallback(callbackQuery *tgbotapi.CallbackQuery, userAppointments *entites.Appointment) {

	deleteMsg := tgbotapi.NewDeleteMessage(callbackQuery.From.ID, callbackQuery.Message.MessageID)
	h.api.Send(deleteMsg)

	file, err := os.Open("../../assets/Коррекция.jpg")
	if err != nil {
		log.Println(err)

		errText := "Произошла ошибка, попробуйте снова позже"
		errMsg := tgbotapi.NewMessage(callbackQuery.From.ID, errText)
		h.api.Send(errMsg)
		return
	}
	defer file.Close()

	reader := tgbotapi.FileReader{Name: "Коррекция.jpg", Reader: file}

	photo := tgbotapi.NewPhoto(callbackQuery.From.ID, reader)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userAppointments.Services, err = h.service.AppointmentService.GetAvailableServices(ctx)
	if err != nil {
		h.api.Send(tgbotapi.NewMessage(callbackQuery.From.ID, "Произошла ошибка, попробуйте позже"))
		log.Println("Ошибка получения услуг: ", err)
		return
	}
	editText := fmt.Sprintf("≪━─◈  "+userAppointments.Services[0].Name+"  ◈─━≫\n\n"+userAppointments.Services[0].Descr+"\nЗанимает %d минут\n\n", userAppointments.Services[0].Duration)
	for i, serv := range userAppointments.Services {
		if i == 0 {
			editText = editText + "☑️	__<u>" + serv.Name + "</u>__	☑️\n"
		} else {
			editText = editText + "☑️	" + serv.Name + "	☑️\n"
		}
	}

	backwardIndex := len(userAppointments.Services) - 1
	forwardIndex := 1

	forwardService := "service_" + strconv.Itoa(forwardIndex)
	backwardService := "service_" + strconv.Itoa(backwardIndex)

	editKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", backwardService),
			tgbotapi.NewInlineKeyboardButtonData("Вперед ➡️", forwardService),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Добавить ✅", "addRemove_0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к времени", "backToTime"),
		),
	)

	photo.ReplyMarkup = editKeyboard
	photo.Caption = editText
	photo.ParseMode = "HTML"

	h.api.Send(photo)

}

func (h *BotHandler) handleServiceChooseCallback(callbackQuery *tgbotapi.CallbackQuery, serviceIndex int, userAppointments *entites.Appointment) {

	deleteMsg := tgbotapi.NewDeleteMessage(callbackQuery.From.ID, callbackQuery.Message.MessageID)
	h.api.Send(deleteMsg)

	editText := fmt.Sprintf("≪━─◈  "+userAppointments.Services[serviceIndex].Name+"  ◈─━≫\n\n"+userAppointments.Services[serviceIndex].Descr+"\nЗанимает %d минут\n\n", userAppointments.Services[serviceIndex].Duration)
	for i, serv := range userAppointments.Services {

		if i == serviceIndex {
			if userAppointments.Services[i].Added {
				editText = editText + "✅	__<u>" + serv.Name + "</u>__	✅\n"
			} else {
				editText = editText + "☑️	__<u>" + serv.Name + "</u>__	☑️\n"
			}
		} else if userAppointments.Services[i].Added {
			editText = editText + "✅	" + serv.Name + "	✅\n"
		} else {
			editText = editText + "☑️	" + serv.Name + "	☑️\n"
		}

	}

	var backwardIndex, forwardIndex int

	switch serviceIndex {
	case len(userAppointments.Services) - 1:
		backwardIndex = serviceIndex - 1
		forwardIndex = 0
	case 0:
		backwardIndex = len(userAppointments.Services) - 1
		forwardIndex = serviceIndex + 1
	default:
		backwardIndex = serviceIndex - 1
		forwardIndex = serviceIndex + 1
	}

	forwardService := "service_" + strconv.Itoa(forwardIndex)
	backwardService := "service_" + strconv.Itoa(backwardIndex)

	addRemoveText := "✅ Добавить ✅"
	if userAppointments.Services[serviceIndex].Added {
		addRemoveText = "❌ Отменить ❌"
	}

	addRemoveData := "addRemove_" + strconv.Itoa(serviceIndex)

	ableToConfirm := false
	for i := range userAppointments.Services {
		if userAppointments.Services[i].Added {
			ableToConfirm = true
			break
		}
	}

	var editKeyboard tgbotapi.InlineKeyboardMarkup

	if ableToConfirm {
		editKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", backwardService),
				tgbotapi.NewInlineKeyboardButtonData("Вперед ➡️", forwardService),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(addRemoveText, addRemoveData),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📅 Записаться 📅", "confirmServices"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к времени", "backToTime"),
			),
		)
	} else {
		editKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", backwardService),
				tgbotapi.NewInlineKeyboardButtonData("Вперед ➡️", forwardService),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(addRemoveText, addRemoveData),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к времени", "backToTime"),
			),
		)
	}

	filePath := "../../assets/" + userAppointments.Services[serviceIndex].Name + ".jpg"
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)

		errText := "Произошла ошибка, попробуйте снова позже"
		errMsg := tgbotapi.NewMessage(callbackQuery.From.ID, errText)
		h.api.Send(errMsg)
		return
	}
	defer file.Close()

	reader := tgbotapi.FileReader{Name: userAppointments.Services[serviceIndex].Name + ".jpg", Reader: file}

	photo := tgbotapi.NewPhoto(callbackQuery.From.ID, reader)

	photo.Caption = editText
	photo.ReplyMarkup = editKeyboard

	photo.ParseMode = "HTML"

	h.api.Send(photo)

}

func (h *BotHandler) handleAddRemoveServiceCallback(callbackQuery *tgbotapi.CallbackQuery, serviceIndex int, userAppointments *entites.Appointment) {
	editText := fmt.Sprintf("≪━─◈  "+userAppointments.Services[serviceIndex].Name+"  ◈─━≫\n\n"+userAppointments.Services[serviceIndex].Descr+"\nЗанимает %d минут\n\n", userAppointments.Services[serviceIndex].Duration)
	for i, serv := range userAppointments.Services {

		if i == serviceIndex {
			if userAppointments.Services[i].Added {
				editText = editText + "✅	__<u>" + serv.Name + "</u>__	✅\n"
			} else {
				editText = editText + "☑️	__<u>" + serv.Name + "</u>__	☑️\n"
			}
		} else if userAppointments.Services[i].Added {
			editText = editText + "✅	" + serv.Name + "	✅\n"
		} else {
			editText = editText + "☑️	" + serv.Name + "	☑️\n"
		}

	}

	var backwardIndex, forwardIndex int

	switch serviceIndex {
	case len(userAppointments.Services) - 1:
		backwardIndex = serviceIndex - 1
		forwardIndex = 0
	case 0:
		backwardIndex = len(userAppointments.Services) - 1
		forwardIndex = serviceIndex + 1
	default:
		backwardIndex = serviceIndex - 1
		forwardIndex = serviceIndex + 1
	}

	forwardService := "service_" + strconv.Itoa(forwardIndex)
	backwardService := "service_" + strconv.Itoa(backwardIndex)

	addRemoveText := "✅ Добавить ✅"
	if userAppointments.Services[serviceIndex].Added {
		addRemoveText = "❌ Отменить ❌"
	}

	addRemoveData := "addRemove_" + strconv.Itoa(serviceIndex)

	ableToConfirm := false
	for i := range userAppointments.Services {
		if userAppointments.Services[i].Added {
			ableToConfirm = true
			break
		}
	}

	var editKeyboard tgbotapi.InlineKeyboardMarkup

	if ableToConfirm {
		editKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", backwardService),
				tgbotapi.NewInlineKeyboardButtonData("Вперед ➡️", forwardService),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(addRemoveText, addRemoveData),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📅 Записаться 📅", "confirmServices"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к времени", "backToTime"),
			),
		)
	} else {
		editKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", backwardService),
				tgbotapi.NewInlineKeyboardButtonData("Вперед ➡️", forwardService),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(addRemoveText, addRemoveData),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к времени", "backToTime"),
			),
		)
	}

	editMsg := tgbotapi.NewEditMessageCaption(callbackQuery.From.ID, callbackQuery.Message.MessageID, editText)
	editMsg.ReplyMarkup = &editKeyboard
	editMsg.ParseMode = "HTML"
	h.api.Send(editMsg)
}

func (h *BotHandler) handleBackToTimeCallback(callbackQuery *tgbotapi.CallbackQuery, userAppointment *entites.Appointment) {
	deleteMsg := tgbotapi.NewDeleteMessage(callbackQuery.From.ID, callbackQuery.Message.MessageID)
	h.api.Send(deleteMsg)

	text := "**" + rusdate.FormatDayMonth(userAppointment.Date) + "**\n\nНа какое время?\n\n"

	var rows [4][4]tgbotapi.InlineKeyboardButton

	startTime := userAppointment.Date.Add(12 * time.Hour)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			btnText := startTime.Format("15:04")
			callBackText := "chooseTime:" + btnText

			// Сделать после || Проверку на наличие записи в это время у БД

			if time.Now().After(startTime.Add(-30 * time.Minute)) {
				btnText = "❌" + btnText
				callBackText = "inactiveTime"
			} else {
				btnText = "☑️" + btnText
			}

			rows[i][j] = tgbotapi.InlineKeyboardButton{
				Text:         btnText,
				CallbackData: &callBackText,
			}

			startTime = startTime.Add(30 * time.Minute)
		}
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows[0][:], rows[1][:], rows[2][:], rows[3][:],
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(functionalButtons.back, "backToDate")))

	msg := tgbotapi.NewMessage(callbackQuery.From.ID, text)
	msg.ReplyMarkup = keyboard

	msg.ParseMode = "markdown"
	h.api.Send(msg)
}

func (h *BotHandler) handleServicesConfirmCallback(callbackQuery *tgbotapi.CallbackQuery, userAppointment *entites.Appointment) {
	deleteMsg := tgbotapi.NewDeleteMessage(callbackQuery.From.ID, callbackQuery.Message.MessageID)
	h.api.Send(deleteMsg)

	confirmText := fmt.Sprintf("Запись:\n%s | %s:%s\n\nУслуги:\n", rusdate.FormatDayMonth(userAppointment.Date), userAppointment.Hour, userAppointment.Minute)

	totalDuration := 0
	for i := range userAppointment.Services {
		if userAppointment.Services[i].Added {
			confirmText = confirmText + userAppointment.Services[i].Name + "\n"
			totalDuration += userAppointment.Services[i].Duration
		}
	}

	switch {
	case totalDuration < 60:
		confirmText = fmt.Sprintf(confirmText+"\nВремя выполнения:\n%d минут", totalDuration)
	case totalDuration == 60:
		confirmText = confirmText + "\nВремя выполнения:\n1 час"
	case totalDuration > 60:
		hoursDuration := totalDuration / 60
		minutesDuration := totalDuration % 60
		confirmText = fmt.Sprintf(confirmText+"\nВремя выполнения:\n%d час, %d минут", hoursDuration, minutesDuration)
	}

	confirmMsg := tgbotapi.NewMessage(callbackQuery.From.ID, confirmText)

	replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить ✅", "confirmAppointment"),
		),
	)

	confirmMsg.ReplyMarkup = replyKeyboard

	h.api.Send(confirmMsg)
}

func (h *BotHandler) handleAppointmentConfirmCallback(callbackQuery *tgbotapi.CallbackQuery, userAppointment *entites.Appointment) {

}

// Добавить фотографии -> Удалять старое сообщение, создавать новое с прикреплением фото.
// Добавить кнопки confirmAppointment и backToTime

// Сделать выбор услуг 1 сообщением, а после давать выбирать доступное время по totalDuration
