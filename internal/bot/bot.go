package bot

import (
	"github.com/Thomas3246/BrowsMasterManager/internal/handle"
	"github.com/Thomas3246/BrowsMasterManager/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func NewBot(apiToken string) (*Bot, error) {
	botApi, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, err
	}
	return &Bot{api: botApi}, nil
}

func (b *Bot) Start() {

	botService := service.NewBotService()
	botHandler := handle.NewBotHandler(b.api, botService)

	botHandler.Start()

}
