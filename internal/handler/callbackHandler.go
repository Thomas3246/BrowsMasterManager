package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Кнопка подтверждения имени при уточнении
func (h *BotHandler) handleConfirmNameCallback(update *tgbotapi.Update) {
	callbackQuery := update.CallbackQuery
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	h.api.Request(callback)

	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Отлично, так и оставим")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	h.api.Send(msg)
}

// Кнопка смены имени при уточнении
func (h *BotHandler) handleChangeNameCallback(update *tgbotapi.Update) {
	callbackQuery := update.CallbackQuery
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	h.api.Request(callback)

	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Для смены имени напишите команду \"/name ___имя___\". \nНапример: \"/name Лена\"")
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	h.api.Send(msg)
}
