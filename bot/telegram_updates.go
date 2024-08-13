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
)

var commands []string = []string{"/start", "/help", "/register", "/stop", "/resume"}

type TelegramUpdates struct {
	bot  *tgbotapi.BotAPI
	db   *DBOps
	conf config.Config
}

func NewTelegramUpdates(b *tgbotapi.BotAPI, conf config.Config) *TelegramUpdates {
	return &TelegramUpdates{
		bot:  b,
		db:   NewDB(),
		conf: conf,
	}
}

type WebhookRequest struct {
	URL         string `json:"url"`
	Certificate string `json:"certificate"`
}

func (tu *TelegramUpdates) SetWebhookOnStart() {
	token := tu.bot.Token
	url := tu.conf.GetString("TELEGRAM_WEBHOOK_URL")

	// Prepare the API endpoint URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?", token)

	certificatestringFromFile, err := os.ReadFile("../certs/pub.pem")
	if err != nil {
		log.Fatalf("Failed to read certificate file: %v", err)
	}

	body := bytes.NewBuffer([]byte{})
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("url", url)
	part, err := writer.CreateFormFile("certificate", "../certs/pub.pem")
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
	token := p.telegramUpdatesInst.RegisterChat(channelchatid)
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
func (tu *TelegramUpdates) RegisterChat(channelid int64) string {
	token := generateRandomTokenWithLength(24)
	// tu.db.SaveChat(channelid, token)
	return token

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
