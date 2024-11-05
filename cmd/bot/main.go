package main

import (
	"log"

	"github.com/Thomas3246/BrowsMasterManager/internal/bot"
	"github.com/Thomas3246/BrowsMasterManager/internal/config"
	"github.com/Thomas3246/BrowsMasterManager/internal/repository/sqlite"
)

func main() {

	// To-DO посмотреть про контексты и мб добавить их

	// Добавить переменное окружение

	// apiToken := os.Getenv("TELEGRAM_API_TOKEN")
	// if apiToken == "" {
	// 	panic("Необходимо установить TELEGRAM_API_TOKEN")
	// }

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
