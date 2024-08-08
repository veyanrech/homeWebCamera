package main

import (
	"os"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/veyanrech/homeWebCamera/imagecapture/config"
)

func main() {

	conf := config.NewConfig()

	// set webhook for bot

	bot, err := tgbot.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		panic(err)
	}

	//listen for updates
}
