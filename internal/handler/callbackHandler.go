package handler

import (
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
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date), "todayDate + 0"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24)), "todayDate + 1"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*2)), "todayDate + 2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*3)), "todayDate + 3"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*4)), "todayDate + 4"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*5)), "todayDate + 5"),
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

	userAppointments.DateTime = date.Add(time.Hour * 24 * time.Duration(dayNumber))

	h.api.Send(editMsg)

}

func (h *BotHandler) handleDateConfirmCallback(callbackQuery *tgbotapi.CallbackQuery, userAppointments *entites.Appointment) {
	editText := rusdate.FormatDayMonth(userAppointments.DateTime) + "\n\nНа какое время?\n\n"

	var rows [4][4]tgbotapi.InlineKeyboardButton

	// startTime := 12 * time.Hour
	startTime := userAppointments.DateTime.Add(12 * time.Hour)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			btnText := startTime.Format("15:04")
			callBackText := "chooseTime:" + btnText

			// Сделать после || Проверку на наличие записи в это время у БД

			if time.Now().After(startTime.Add(-30 * time.Minute)) {
				btnText = "~~" + btnText + "~~"
				callBackText = "inactiveTime"
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

	h.api.Send(editMsg)
}

func (h *BotHandler) handleBackToDate(callbackQuery *tgbotapi.CallbackQuery) {
	msgText := "Давайте выберем дату\n\n"

	date := time.Now()
	dateKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			// tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(date.Day()) + date., ""),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date), "todayDate + 0"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24)), "todayDate + 1"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*2)), "todayDate + 2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*3)), "todayDate + 3"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*4)), "todayDate + 4"),
			tgbotapi.NewInlineKeyboardButtonData(rusdate.FormatDayMonth(date.Add(time.Hour*24*5)), "todayDate + 5"),
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

func (h *BotHandler) handleTimeChooseCallback(callbackQuery *tgbotapi.CallbackQuery, hour int, minute int, userAppointments *entites.Appointment) {

	// edit like date choose
}
