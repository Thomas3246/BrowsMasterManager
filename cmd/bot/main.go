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
	if cfg.MasterPhone == "" {
		log.Fatalf("Номер мастера не введен")
	}
	if cfg.AboutMaster == "" {
		log.Fatalf("Информация о мастере не введена")
	}

	db, err := sqlite.InitDB()
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	// Добавить в NewBot и NewService с параметром db параметр redis

	bot, err := bot.NewBot(cfg, db)
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	err = bot.Start()
	if err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}
}
