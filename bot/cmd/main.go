package main

import (
	"fmt"
	"net/http"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/veyanrech/homeWebCamera/bot"
	"github.com/veyanrech/homeWebCamera/imagecapture/config"
)

func main() {

	conf := config.NewConfig()

	if conf == nil { //TODO: on depluy should change to another way of getting confic
		panic("Config not found")
	}

	// send request to telegram to set a webgook

	_bot, err := tgbot.NewBotAPI(conf.GetString("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	botInst := bot.NewTelegramUpdates(_bot, conf)
	botInst.SetWebhookOnStart()
	updchan := _bot.ListenForWebhook("/webhook")

	go func() {
		for {
			upd := <-updchan
			botInst.ProcessUpdates(upd)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServeTLS(fmt.Sprintf(":%s", conf.GetString("BOT_PORT")), "../certs/pub.pem", "../certs/private.key", nil)

	select {}

}
