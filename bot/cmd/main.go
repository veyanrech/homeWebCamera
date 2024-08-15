package main

import (
	"fmt"
	"net/http"
	"os"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/veyanrech/homeWebCamera/bot"
	"github.com/veyanrech/homeWebCamera/imagecapture/config"
	"github.com/veyanrech/homeWebCamera/imagecapture/utils"
)

func main() {

	var conf config.Config
	var filename string

	osgetwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	switch opsys := utils.GetOS(); opsys {
	default:
		filename = osgetwd + "/macos.config.json"
	case "windows":
		filename = osgetwd + string(os.PathSeparator) + "win.config.json"
	}

	conf = config.NewConfig(filename)

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

	app := bot.NewFilesReceriverClient(botInst)

	http.HandleFunc("/filesreciever", app.RecieveFileHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		respmap := make(map[string]string)

		//check db conenction also
		err := botInst.PingDB()
		if err != nil {
			respmap["db"] = "error"
		} else {
			respmap["db"] = "ok"
		}

		respmap["bot"] = "ok"

		w.WriteHeader(http.StatusOK)

		for k, v := range respmap {
			w.Write([]byte(fmt.Sprintf("%s: %s\n", k, v)))
		}

	})

	go func() {
		for {
			upd := <-updchan
			botInst.ProcessUpdates(upd)
		}
	}()

	//certs are in the same folder as the binary
	publocation := osgetwd + "/certs/pub.pem"
	privlocation := osgetwd + "/certs/private.key"

	http.ListenAndServeTLS(fmt.Sprintf(":%s", conf.GetString("BOT_PORT")), publocation, privlocation, nil)

	select {}

}
