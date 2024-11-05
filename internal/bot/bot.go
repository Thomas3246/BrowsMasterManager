package bot

import (
	"database/sql"

	"github.com/Thomas3246/BrowsMasterManager/internal/handler"
	"github.com/Thomas3246/BrowsMasterManager/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
	db  *sql.DB
}

func NewBot(apiToken string, db *sql.DB) (*Bot, error) {
	botApi, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, err
	}
	return &Bot{
		api: botApi,
		db:  db,
	}, nil
}

func (b *Bot) Start() {

	botService := service.NewBotService(b.db)
	botHandler := handler.NewBotHandler(b.api, botService)

	botHandler.Start()

}
