package handler

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *BotHandler) AuthMiddleWare(next func(ctx context.Context, update *tgbotapi.Update)) func(ctx context.Context, update *tgbotapi.Update) {
	return func(ctx context.Context, update *tgbotapi.Update) {
		_, isRegistred, err := h.CheckForUser(ctx, update)
		if err != nil {
			log.Println(err)

			reply := tgbotapi.NewMessage(update.FromChat().ID, "Произошла ошибка, попробуйте позже")
			h.api.Send(reply)
			return
		}

		if !isRegistred {
			reply := tgbotapi.NewMessage(update.FromChat().ID, "Эта опция доступна только для авторизованных пользователей, подтвердивших свой номер телефона.\n\nДля этого нажмите на кнопку ниже и поделитесь своим номером телефона")
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButtonContact("Поделиться"),
				),
			)
			reply.ReplyMarkup = keyboard
			h.api.Send(reply)
			return
		}

		next(ctx, update)
	}
}
