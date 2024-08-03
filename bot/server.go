package bot

import "net/http"

type BotServer struct {
}

func NewBotServer() *BotServer {
	return &BotServer{}
}

func (s *BotServer) Start() {}

func (s *BotServer) Stop() {}

func (s *BotServer) HandleChatRegistration(req *http.Request, res http.ResponseWriter) {

}
