package bot

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/veyanrech/homeWebCamera/imagecapture/config"
	"github.com/veyanrech/homeWebCamera/imagecapture/utils"
)

var commands []string = []string{"/start", "/help", "/register", "/stop", "/resume"}

type TelegramUpdates struct {
	log  utils.Logger
	bot  *tgbotapi.BotAPI
	db   *DBOps
	conf config.Config
}

func NewTelegramUpdates(b *tgbotapi.BotAPI, conf config.Config) *TelegramUpdates {
	res := &TelegramUpdates{
		bot:  b,
		db:   NewDB(conf),
		conf: conf,
		log:  utils.NewFileLogger("telegram_updates.log", "telegram_updates_errors.log"),
	}

	if os.Getenv("ENVIRONMENT") == "PROD" {
		res.log.Disable()
	}

	return res
}

type WebhookRequest struct {
	URL         string `json:"url"`
	Certificate string `json:"certificate"`
}

func (tu *TelegramUpdates) PingDB() error {
	return tu.db.Ping()
}

func (tu *TelegramUpdates) SetWebhookOnStart() {
	token := tu.bot.Token
	url := tu.conf.GetString("TELEGRAM_WEBHOOK_URL")

	osgetwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Prepare the API endpoint URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?", token)

	certificatestringFromFile, err := os.ReadFile(osgetwd + "/certs/pub.pem")
	if err != nil {
		log.Fatalf("Failed to read certificate file: %v", err)
	}

	body := bytes.NewBuffer([]byte{})
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("url", url)
	part, err := writer.CreateFormFile("certificate", osgetwd+"/certs/pub.pem")
	if err != nil {
		log.Fatalf("Failed to create form file: %v", err)
	}

	_, err = io.Copy(part, bufio.NewReader(bytes.NewBuffer(certificatestringFromFile)))
	if err != nil {
		log.Fatalf("Failed to copy file: %v", err)
	}

	err = writer.Close()
	if err != nil {
		log.Fatalf("Failed to close writer: %v", err)
	}

	// Create a new HTTP request
	//send the certificate

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Set the content type to JSON
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to set webhook: %v", resp.Status)
	}
}

func (tu *TelegramUpdates) ProcessUpdates(u tgbotapi.Update) {

	//bot waits for messages from users or chat

	mesproc := tu.newMessageProcessor(&u)

	mesproc.processMessage()
}

type messageProcessorI interface {
	processMessage()
}

func (tu *TelegramUpdates) newMessageProcessor(u *tgbotapi.Update) messageProcessorI {
	if u.Message != nil {
		return &privateMessageProcessor{updateInfo: u, telegramUpdatesInst: tu}
	} else {
		if u.ChannelPost != nil {
			return &groupMessageProcessor{updateInfo: u, telegramUpdatesInst: tu}
		}
	}

	return nil
}

type privateMessageProcessor struct {
	updateInfo          *tgbotapi.Update
	telegramUpdatesInst *TelegramUpdates
}

func (p *privateMessageProcessor) processMessage() {

	m := p.updateInfo.Message

	if m == nil {
		return
	}

	if m.IsCommand() {
		switch m.Command() {
		case "start", "help":
			p.telegramUpdatesInst.SendCommandsList(m.Chat.ID)
		case "stop":
			p.telegramUpdatesInst.Stop()
		case "resume":
			p.telegramUpdatesInst.Resume()
		case "register":
			p.registerChat()
		}
	}
}

func (p *privateMessageProcessor) registerChat() {
	channelchatid := p.updateInfo.Message.ForwardFromChat.ID
	token, err := p.telegramUpdatesInst.RegisterChat(channelchatid)
	if err != nil {
		p.telegramUpdatesInst.db.logger.Error(fmt.Sprintf("Error while registering chat: %v", err))
		p.telegramUpdatesInst.bot.Send(tgbotapi.NewMessage(p.updateInfo.Message.Chat.ID, "Error while registering chat"))
		return
	}
	chatid := p.updateInfo.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatid, fmt.Sprintf("Your token is: %s", token))
	p.telegramUpdatesInst.bot.Send(msg)
}

type groupMessageProcessor struct {
	updateInfo          *tgbotapi.Update
	telegramUpdatesInst *TelegramUpdates
}

func (p *groupMessageProcessor) processMessage() {

}

func (tu *TelegramUpdates) SendCommandsList(chatid int64) {
	var commandsList string
	for _, c := range commands {
		commandsList += c + "\n"
	}

	msg := tgbotapi.NewMessage(chatid, commandsList)
	tu.bot.Send(msg)
}

// RegisterChat register chat to send photos
// this method will return the randomly generated token
func (tu *TelegramUpdates) RegisterChat(channelid int64) (string, error) {
	//check if chat is already registered
	//if registered return the token
	//if not registered generate a token and register the chat
	res, err := tu.db.FindChatID(channelid)
	if err != nil {
		tu.db.logger.Error(fmt.Sprintf("Error while finding chat id: %v", err))
		return "", err
	}

	if res.token != "" {
		return res.token, nil
	}

	token := generateRandomTokenWithLength(24)
	tu.db.RegisterChatID(channelid, token)
	return token, nil

}

func (tu *TelegramUpdates) SendFileToChat(fileinput map[string][]*multipart.FileHeader, chatidtosend int64) {
	mediagroupConfig := tgbotapi.NewMediaGroup(chatidtosend, []interface{}{})

	filestosend := []interface{}{}

	for k, fileheader := range fileinput {
		tu.log.Info(fmt.Sprintf("file: %s", k))
		for _, fileh := range fileheader {
			tu.log.Info(fmt.Sprintf("file: %s", fileh.Filename))
			file, err := fileh.Open()
			if err != nil {
				tu.log.Error(fmt.Sprintf("Error while opening file: %v", err))
				return
			}
			defer file.Close()

			fileBytes := bytes.Buffer{}
			_, err = io.Copy(&fileBytes, file)
			if err != nil {
				tu.log.Error(fmt.Sprintf("Error while reading file: %v", err))
				return
			}

			tgfilebytes := tgbotapi.FileBytes{
				Name:  fileh.Filename,
				Bytes: fileBytes.Bytes(),
			}

			tgbotfile := tgbotapi.NewInputMediaPhoto(tgfilebytes)

			filestosend = append(filestosend, tgbotfile)
		}
	}

	mediagroupConfig.Media = filestosend

	//send file to chat
	_, err := tu.bot.SendMediaGroup(mediagroupConfig)
	if err != nil {
		tu.log.Error(fmt.Sprintf("Error while sending media group: %v", err))
		return
	}
}

// stop sending photos to chat
func (tu *TelegramUpdates) Stop() {}

// resume sending photos to chat
func (tu *TelegramUpdates) Resume() {}

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomTokenWithLength(length int) string {
	randomStringResult := make([]rune, length)
	for i := 0; i < length; i++ {
		randomStringResult[i] = letters[rand.Intn(len(letters))]
	}

	return string(randomStringResult)
}
