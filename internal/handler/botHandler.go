package handler

import (
	"log"
	"strconv"
	"strings"

	"github.com/Thomas3246/BrowsMasterManager/internal/entites"
	"github.com/Thomas3246/BrowsMasterManager/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	api     *tgbotapi.BotAPI
	service *service.BotService
}

type mainButtonStrings struct {
	newAppointment    string
	myAppointments    string
	cancelAppointment string
	serviceList       string
	aboutMaster       string
	help              string
	confirm           string
	back              string
}

var functionalButtons = mainButtonStrings{
	newAppointment:    "📅 Записаться",
	myAppointments:    "📝 Мои записи",
	cancelAppointment: "❌ Отменить запись",
	serviceList:       "💆‍♀️ Услуги",
	aboutMaster:       "💁‍♀ О мастере",
	help:              "❓ Помощь",
	confirm:           "✅ Подтвердить ✅",
	back:              "⬅️ Назад",
}

var usersAppointments = make(map[int64]*entites.Appointment)

func NewBotHandler(api *tgbotapi.BotAPI, service *service.BotService) *BotHandler {
	return &BotHandler{
		api:     api,
		service: service}
}

func (h *BotHandler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := h.api.GetUpdatesChan(u)

	for update := range updates {
		// Для каждого сообщения создается своя горутина для параллелизма
		go h.HandleMessage(&update)
	}
}

func (h *BotHandler) HandleMessage(update *tgbotapi.Update) {

	// Обрабатывается отправка пользователем контакта
	if update.Message != nil {
		if update.Message.Contact != nil {
			h.handleContact(update)
		} else if update.Message.IsCommand() {

			// Обрабаывается отправка пользователем команд
			switch update.Message.Command() {

			case "start":
				// MiddleWare не требуется, общедоступная команда
				h.handleStartCommand(update)

			case "appointment":
				// Для выполнения записи сначала производится проверка, зарегистрирован ли пользователь,
				// после чего проверяется, ввел ли пользователь свое имя
				usersAppointments[update.FromChat().ID] = &entites.Appointment{}
				handler := h.AuthMiddleWare(h.NameMiddleWare(h.handleNewAppointmentCommand))
				handler(update)

			case "name":
				// Для изменения имени необходимо подтверждение, что пользователь зарегистрирован
				handler := h.AuthMiddleWare(h.handleNameChangeCommand)
				handler(update)
			}

		} else {
			switch update.Message.Text {
			case functionalButtons.newAppointment:
				usersAppointments[update.FromChat().ID] = &entites.Appointment{}
				handler := h.AuthMiddleWare(h.NameMiddleWare(h.handleNewAppointmentCommand))
				handler(update)
			}
		}
	}

	// Обработка нажатий на кнопки
	if update.CallbackQuery != nil {
		callbackQuery := update.CallbackQuery
		switch {
		case callbackQuery.Data == "callbackConfirmName":
			h.handleConfirmNameCallback(update)

		case callbackQuery.Data == "callbackChangeName":
			h.handleChangeNameCallback(update)

		case strings.HasPrefix(callbackQuery.Data, "todayDate"):
			parts := strings.Split(callbackQuery.Data, " + ")
			dayNumber, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Println("Ошибка при разборе: ", err)
				break
			}

			h.handleDateChooseCallback(callbackQuery, dayNumber, usersAppointments[callbackQuery.From.ID])

		case callbackQuery.Data == "confirmDate":
			h.handleDateConfirmCallback(callbackQuery, usersAppointments[callbackQuery.From.ID])

		case callbackQuery.Data == "backToDate":
			h.handleBackToDate(callbackQuery)

		case strings.HasPrefix(callbackQuery.Data, "chooseTime"):
			parts := strings.Split(callbackQuery.Data, ":")
			hour, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Println("Ошибка при разборе: ", err)
				break
			}
			minute, err := strconv.Atoi(parts[2])
			if err != nil {
				log.Println("Ошибка при разборе: ", err)
				break
			}

			h.handleTimeChooseCallback(callbackQuery, hour, minute, usersAppointments[callbackQuery.From.ID])
		}

	}
}

func attachFunctionalButtons(msg *tgbotapi.MessageConfig) {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(functionalButtons.newAppointment)),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(functionalButtons.myAppointments),
			tgbotapi.NewKeyboardButton(functionalButtons.cancelAppointment),
		),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(functionalButtons.aboutMaster)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(functionalButtons.help)),
	)
	msg.ReplyMarkup = keyboard
}
