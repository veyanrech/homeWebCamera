package main

import (
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/veyanrech/homeWebCamera/imagecapture/config"
)

func main() {

	conf := config.NewConfig()

	// set webhook for bot

	bot, err := tgbot.NewBotAPI(conf.GetString("telegram_bot_token"))
	if err != nil {
		panic(err)
	}

	//listen for updates
}
