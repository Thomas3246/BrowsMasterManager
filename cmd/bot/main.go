package main

import (
	"log"

	"github.com/Thomas3246/BrowsMasterManager/internal/bot"
)

func main() {
	// apiToken := os.Getenv("TELEGRAM_API_TOKEN")
	// if apiToken == "" {
	// 	panic("Необходимо установить TELEGRAM_API_TOKEN")
	// }

	apiToken := ""

	bot, err := bot.NewBot(apiToken)
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	bot.Start()
}
