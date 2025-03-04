package bot

import (
	"database/sql"

	"github.com/Thomas3246/BrowsMasterManager/internal/config"
	"github.com/Thomas3246/BrowsMasterManager/internal/handler"
	"github.com/Thomas3246/BrowsMasterManager/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

type Bot struct {
	api         *tgbotapi.BotAPI
	db          *sql.DB
	redisClient *redis.Client
	cfg         *config.Config
}

func NewBot(cfg *config.Config, db *sql.DB) (*Bot, error) {
	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	return &Bot{
		api: botApi,
		db:  db,
		cfg: cfg,
	}, nil
}

func (b *Bot) Start() (err error) {

	b.redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	botService := service.NewBotService(b.db, b.redisClient)

	err = botService.UserService.SetMasterPhone(b.cfg.MasterPhone)
	if err != nil {
		return err
	}

	botHandler := handler.NewBotHandler(b.api, botService)

	botHandler.Start()

	return nil
}
