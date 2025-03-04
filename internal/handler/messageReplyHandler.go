package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Горутины и контексты для выполнения самой логики и его ограничения (если задержка именно в выполнении логики а не в запросах к БД,
// так как если запрос к БД будет слишком долгим, он сам прервется и вернет ошибку.)

// (больше как демонстрация, поскольку никаких особо сложных операций не производится)

// Обработчик команды /start
func (h *BotHandler) handleStartCommand(update *tgbotapi.Update) {
	message := update.Message

	startMsg := `Привет! Это бот для записи на брови к мастеру ___ИмяМастера___ по адресу **г.Город, ул.**

		Для записи необходимо подтверждение номера. Для подтверждения нажмите на кнопку *"Поделиться"*.`

	reply := tgbotapi.NewMessage(message.Chat.ID, startMsg)
	reply.ParseMode = "markdown"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("Поделиться"),
		),
	)
	reply.ReplyMarkup = keyboard
	h.api.Send(reply)

	// После регистрации должно выходить сообщение с инструкциями к пользованию ботом и командами. Должна быть команда /help (она и должна вызываться).

}

// Обработчик отправки пользователем контакта
func (h *BotHandler) handleContact(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message := update.Message
	contact := message.Contact

	// Если был отправлен не свой контакт, фукнция прерывается
	if contact.UserID != message.From.ID {
		msg := "Пожалуйста, отправьте именно ваш контакт, используя кнопку \"**Поделиться**\"."
		reply := tgbotapi.NewMessage(message.Chat.ID, msg)
		reply.ParseMode = "markdown"
		h.api.Send(reply)
		return
	}

	// Канал для получения результата горутины
	resultChan := make(chan struct {
		userName    string
		isRegistred bool
		err         error
	})

	go func() {
		userName, isRegistred, err := h.CheckForUser(ctx, update)
		resultChan <- struct {
			userName    string
			isRegistred bool
			err         error
		}{userName, isRegistred, err}
	}()

	select {
	case <-ctx.Done():
		reply := tgbotapi.NewMessage(message.Chat.ID, "Не удалось обработать\nПревышено время ожидания")
		h.api.Send(reply)
		return
	case result := <-resultChan:
		if result.err != nil {
			reply := tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка, попробуйте позже")
			h.api.Send(reply)
			return
		}

		if result.isRegistred {
			if result.userName != "" {
				msg := "Вы уже зарегистрированы, вас зовут " + result.userName + ", верно?"

				reply := tgbotapi.NewMessage(message.Chat.ID, msg)
				reply.ParseMode = "markdown"

				inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Да", "callbackConfirmName"),
						tgbotapi.NewInlineKeyboardButtonData("Нет, изменить", "callbackChangeName"),
					),
				)
				reply.ReplyMarkup = inlineKeyboard

				h.api.Send(reply)
				return
			}

			reply := tgbotapi.NewMessage(message.Chat.ID, "Вы зарегистрированы, но у Вас не указано имя. \nУкажите его через \"/name ___имя___\"")
			reply.ParseMode = "markdown"
			reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			h.api.Send(reply)
			return
		}

		h.handleRegister(ctx, contact)
		return
	}
}

// Функция для обработки регистрации
func (h *BotHandler) handleRegister(ctx context.Context, contact *tgbotapi.Contact) {
	resultChan := make(chan string, 1)

	go func() {
		result := h.registerUser(ctx, contact)
		resultChan <- result
	}()

	select {
	case <-ctx.Done():
		reply := tgbotapi.NewMessage(contact.UserID, "Не удалось обработать\nПревышено время ожидания")
		h.api.Send(reply)
		return
	case result := <-resultChan:
		reply := tgbotapi.NewMessage(contact.UserID, result)
		reply.ParseMode = "markdown"
		reply.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		h.api.Send(reply)
		return
	}
}

func (h *BotHandler) handleNewAppointmentCommand(update *tgbotapi.Update) {

	file, err := os.Open("../../assets/Коррекция.jpg")
	if err != nil {
		log.Println(err)

		errText := "Произошла ошибка, попробуйте снова позже"
		errMsg := tgbotapi.NewMessage(update.FromChat().ID, errText)
		h.api.Send(errMsg)
		return
	}
	defer file.Close()

	reader := tgbotapi.FileReader{Name: "Коррекция.jpg", Reader: file}

	photo := tgbotapi.NewPhoto(update.FromChat().ID, reader)

	services := usersAppointments[update.FromChat().ID].Services
	text := fmt.Sprintf("≪━─◈  "+services[0].Name+"  ◈─━≫\n\n"+services[0].Descr+"\nЗанимает %d минут\n\n", services[0].Duration)
	for i, serv := range services {
		if i == 0 {
			text = text + "☑️	__<u>" + serv.Name + "</u>__	☑️\n"
		} else {
			text = text + "☑️	" + serv.Name + "	☑️\n"
		}
	}

	backwardIndex := len(services) - 1
	forwardIndex := 1

	forwardService := "service_" + strconv.Itoa(forwardIndex)
	backwardService := "service_" + strconv.Itoa(backwardIndex)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", backwardService),
			tgbotapi.NewInlineKeyboardButtonData("Вперед ➡️", forwardService),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Добавить ✅", "addRemove_0"),
		),
	)

	photo.ReplyMarkup = keyboard
	photo.Caption = text
	photo.ParseMode = "HTML"

	h.api.Send(photo)
}

// Обработчик команды /name
func (h *BotHandler) handleNameChangeCommand(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message := update.Message

	// Если была отправлена команда без аргументов
	args := message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Для смены имени необходимо после команды /name указать свое имя. Например: \"/name Лена\"")
		h.api.Send(msg)
		return
	}

	resultChan := make(chan string, 1)
	go func() {
		result := h.changeUserName(ctx, message)
		resultChan <- result
	}()

	select {
	case <-ctx.Done():
		reply := tgbotapi.NewMessage(message.Chat.ID, "Не удалось изменить имя.\nПревышено время ожидания")
		h.api.Send(reply)
		return

	case result := <-resultChan:
		reply := tgbotapi.NewMessage(message.Chat.ID, result)
		h.api.Send(reply)
	}

	useReply := tgbotapi.NewMessage(message.Chat.ID, "Для записи или просмотра информации о услугах, мастере или боте воспользуйтесь кнопками в клавиатуре")
	attachFunctionalButtons(&useReply)
	h.api.Send(useReply)
}

func (h *BotHandler) handleDiscardAppointmentCommand(update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	appointments, err := h.service.UserService.CheckForAppointments(ctx, update.Message.From.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			noRowsMsg := tgbotapi.NewMessage(update.Message.From.ID, "У вас нет активных записей")
			h.api.Send(noRowsMsg)
			return
		}
		log.Println("Ошибка получения записей: ", err)
		errMsg := tgbotapi.NewMessage(update.Message.From.ID, "Произошла ошибка, повторите позже")
		h.api.Send(errMsg)
		return
	}

	if len(appointments) == 0 {
		noAppointmentsMsg := tgbotapi.NewMessage(update.Message.From.ID, "У вас нет активных записей")
		h.api.Send(noAppointmentsMsg)
		return
	}

	for i := range appointments {
		appointments[i].Services, err = h.service.ServiceService.GetServicesInAppointment(ctx, appointments[i].ID)
		if err != nil {
			log.Printf("Ошибка определения услуг на запись: %v", err)
			errMsg := tgbotapi.NewMessage(update.Message.From.ID, "Произошла ошибка, повторите позже")
			h.api.Send(errMsg)
			return
		}
	}

	err = h.service.AppointmentService.SetAppointmentsInCash(ctx, update.Message.From.ID, appointments)
	if err != nil {
		errMsg := tgbotapi.NewMessage(update.Message.From.ID, "Произошла ошибка, повторите позже")
		h.api.Send(errMsg)
		return
	}

	msgText := fmt.Sprintf("Запись на\n%s\n%s:%s\n\nУслуги:\n", appointments[0].DateStr, appointments[0].Hour, appointments[0].Minute)
	for _, service := range appointments[0].Services {
		msgText = msgText + service.Name + "\n"
	}
	msgText = fmt.Sprintf(msgText+"\nДлительность: %d минут\nСтоимость: %d ₽", appointments[0].TotalDuration, appointments[0].TotalCost)

	msg := tgbotapi.NewMessage(update.Message.From.ID, msgText)

	cancelCallbackText := fmt.Sprintf("confirmCancelAppointment_%d", appointments[0].ID)
	if len(appointments) == 1 {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌ Отменить ❌", cancelCallbackText)))
	} else {
		changeCancelAppointmentText := fmt.Sprintf("changeCancelAppointment_%d", appointments[1].ID)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(" ➡️ ", changeCancelAppointmentText)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌ Отменить ❌", cancelCallbackText)),
		)
	}

	h.api.Send(msg)

}

func (h *BotHandler) handleAboutMasterCommand(update *tgbotapi.Update) {
	aboutMaster, err := h.service.UserService.GetAboutMaster()
	if err != nil {
		errMsg := tgbotapi.NewMessage(update.Message.From.ID, "Произошла ошибка. Пожалуйста, попробуйте позже")
		h.api.Send(errMsg)
	}

	msg := tgbotapi.NewMessage(update.Message.From.ID, aboutMaster)
	h.api.Send(msg)
}

func (h *BotHandler) handleHelpCommand(update *tgbotapi.Update) {
	msgText := "Команда /help - помощь\nКоманда /start - запуск бота. В случае, если нет функциональных кнопок, необходимо перезапустить бота.\n"
	msgText = msgText + "Для того, чтобы подтвердить свой номер телефона, необходимо отправить команду /start, после чего поделиться своим номером телефона"
	msgText = msgText + "Для смены или указания имени необходимо выполнить команду вида \"/name имя\""
	msg := tgbotapi.NewMessage(update.Message.From.ID, msgText)
	h.api.Send(msg)
}
