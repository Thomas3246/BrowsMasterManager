package handler

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

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

func (h *BotHandler) newAppointment(id int64) (*entites.Appointment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	services, err := h.service.AppointmentService.GetAvailableServices(ctx)
	if err != nil {
		log.Println("Ошибка получения услуг: ", err)
		return nil, err
	}
	return &entites.Appointment{Services: services, UserId: id}, nil
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
				var err error
				usersAppointments[update.FromChat().ID], err = h.newAppointment(update.FromChat().ID)
				if err != nil {
					errMsg := tgbotapi.NewMessage(update.FromChat().ID, "Произошла ошибка, попробуйте снова позже")
					h.api.Send(errMsg)
					return
				}
				handler := h.AuthMiddleWare(h.NameMiddleWare(h.handleNewAppointmentCommand))
				handler(update)

			case "name":
				// Для изменения имени необходимо подтверждение, что пользователь зарегистрирован
				handler := h.AuthMiddleWare(h.handleNameChangeCommand)
				handler(update)
			}

		} else {

			// Сюда добавить проверку состояния FSM из map
			// case StateChangeName:
			// 	newName := update.Message.Text
			// 	h.saveNewUserName(userID, newName)
			// 	h.setUserState(userID, StateStart)
			//
			// типа такого

			switch update.Message.Text {
			case functionalButtons.newAppointment:
				var err error
				usersAppointments[update.FromChat().ID], err = h.newAppointment(update.FromChat().ID)
				if err != nil {
					errMsg := tgbotapi.NewMessage(update.FromChat().ID, "Произошла ошибка, попробуйте снова позже")
					h.api.Send(errMsg)
					return
				}
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

		case strings.HasPrefix(callbackQuery.Data, "service_"):
			service_str := strings.TrimPrefix(callbackQuery.Data, "service_")
			service_id, _ := strconv.Atoi(service_str)
			if usersAppointments[callbackQuery.From.ID] != nil {
				h.handleServiceChooseCallback(callbackQuery, service_id, usersAppointments[callbackQuery.From.ID])
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}

		case strings.HasPrefix(callbackQuery.Data, "addRemove_"):
			service_str := strings.TrimPrefix(callbackQuery.Data, "addRemove_")
			service_id, _ := strconv.Atoi(service_str)
			if usersAppointments[callbackQuery.From.ID] != nil {
				usersAppointments[callbackQuery.From.ID].Services[service_id].Added = !usersAppointments[callbackQuery.From.ID].Services[service_id].Added
				h.handleAddRemoveServiceCallback(callbackQuery, service_id, usersAppointments[callbackQuery.From.ID])
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}

		case callbackQuery.Data == "confirmServices":
			if usersAppointments[callbackQuery.From.ID] != nil {
				h.handleServicesConfirmCallback(callbackQuery)
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}

		case strings.HasPrefix(callbackQuery.Data, "date_"):
			dayNumberStr := strings.TrimPrefix(callbackQuery.Data, "date_")
			dayNumber, err := strconv.Atoi(dayNumberStr)
			if err != nil {
				log.Println(err)
			}

			if usersAppointments[callbackQuery.From.ID] != nil {
				h.handleDateChooseCallback(callbackQuery, dayNumber, usersAppointments[callbackQuery.From.ID])
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}

		case callbackQuery.Data == "backToServices":
			if usersAppointments[callbackQuery.From.ID] != nil {
				h.handleBackToServicesCallback(update)
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}

		case callbackQuery.Data == "confirmDate":
			if usersAppointments[callbackQuery.From.ID] != nil {
				h.handleDateConfirmCallback(callbackQuery, usersAppointments[callbackQuery.From.ID])
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}

		case strings.HasPrefix(callbackQuery.Data, "chooseTime"):
			parts := strings.Split(callbackQuery.Data, ":")
			hour := parts[1]
			minute := parts[2]

			if usersAppointments[callbackQuery.From.ID] != nil {
				h.handleTimeChooseCallback(callbackQuery, hour, minute, usersAppointments[callbackQuery.From.ID])
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}

		case callbackQuery.Data == "inactiveTime":
			alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "На данное время запись невозможна\n\nПожалуйста, выберите другое время\n\nДоступное время помечено ☑️")
			h.api.Send(alert)

		case callbackQuery.Data == "backToDate":
			if usersAppointments[callbackQuery.From.ID] != nil {
				h.handleBackToDate(callbackQuery)
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}

		case callbackQuery.Data == "confirmTime":
			if usersAppointments[callbackQuery.From.ID] != nil {
				h.handleTimeConfirmCallback(callbackQuery, usersAppointments[callbackQuery.From.ID])
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}

		case callbackQuery.Data == "backToTime":
			if usersAppointments[callbackQuery.From.ID] != nil {
				h.handleDateConfirmCallback(callbackQuery, usersAppointments[callbackQuery.From.ID])
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}

		case callbackQuery.Data == "confirmAppointment":
			if usersAppointments[callbackQuery.From.ID] != nil {
				h.handleAppointmentConfirmCallback(callbackQuery, usersAppointments[callbackQuery.From.ID])
			} else {
				alert := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Пожалуйста, начните новую запись")
				h.api.Send(alert)
			}
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
