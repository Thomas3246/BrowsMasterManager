package main

import (
	"log"

	"github.com/Thomas3246/BrowsMasterManager/internal/bot"
	"github.com/Thomas3246/BrowsMasterManager/internal/config"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository/sqlite"
)

func main() {

	cfg, err := config.NewConfig("../../config/config.json")
	if err != nil {
		log.Fatalf("Ошибка при инициализации конфиг-файла: %v", err)
	}
	if cfg.BotToken == "" {
		log.Fatalf("API-токен не введён")
	}
	apiToken := cfg.BotToken

	db, err := sqlite.InitDB()
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	bot, err := bot.NewBot(apiToken, db)
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	bot.Start()
}
