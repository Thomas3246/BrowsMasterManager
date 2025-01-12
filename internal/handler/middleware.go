package handler

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MiddleWare для проверки регистрации (регистрация производится при отправке своего контакта боту)
func (h *BotHandler) AuthMiddleWare(next func(update *tgbotapi.Update)) func(update *tgbotapi.Update) {
	return func(update *tgbotapi.Update) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, isRegistred, err := h.CheckForUser(ctx, update)
		if err != nil {
			log.Println(err)

			reply := tgbotapi.NewMessage(update.FromChat().ID, "Произошла ошибка, попробуйте позже")
			h.api.Send(reply)
			return
		}

		if !isRegistred {
			msgText := "Эта опция доступна только для авторизованных пользователей, подтвердивших свой номер телефона.\n\nДля этого нажмите на кнопку ниже и поделитесь своим номером телефона"
			reply := tgbotapi.NewMessage(update.FromChat().ID, msgText)
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButtonContact("Поделиться"),
				),
			)
			reply.ReplyMarkup = keyboard
			h.api.Send(reply)
			return
		}

		next(update)
	}
}

func (h *BotHandler) NameMiddleWare(next func(update *tgbotapi.Update)) func(update *tgbotapi.Update) {
	return func(update *tgbotapi.Update) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userName, _, err := h.CheckForUser(ctx, update)
		if err != nil {
			log.Println(err)
			reply := tgbotapi.NewMessage(update.FromChat().ID, "Произошла ошибка, попробуйте позже")
			h.api.Send(reply)
			return
		}

		if userName == "" {
			msgText := "Для выполнения этой опции необходимо указать свое имя командой *\"/name\"*"
			reply := tgbotapi.NewMessage(update.FromChat().ID, msgText)
			reply.ParseMode = "markdown"
			h.api.Send(reply)
			return
		}

		next(update)
	}
}
