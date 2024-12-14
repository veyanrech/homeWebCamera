package bot

import (
	"fmt"
	"net/http"
	"sync"
)

type FilesReceriverClient struct {
	telegramapiInst *TelegramUpdates
	tokenCache      map[string]int64
	sncmtx          sync.Mutex
}

func NewFilesReceriverClient(t *TelegramUpdates) *FilesReceriverClient {
	return &FilesReceriverClient{
		telegramapiInst: t,
		tokenCache:      make(map[string]int64),
		sncmtx:          sync.Mutex{},
	}
}

func (f *FilesReceriverClient) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//check header
	if r.Header.Get("X-Chat-Registration-Token") == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chatinfo := &chatInfo{}
	var err error

	chatinfo.chatID, err = f.returnChatID(r.Header.Get("X-Chat-Registration-Token"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	data := r.MultipartForm.Value
	if len(data) == 0 {
		http.Error(w, "No data found", http.StatusBadRequest)
		return
	}

	message, ok := data["pcmessage"]
	if !ok {
		http.Error(w, "No message found", http.StatusBadRequest)
		return
	}

	f.telegramapiInst.SendText(chatinfo.chatID, message[0])
}

func (f *FilesReceriverClient) RecieveFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//check header
	if r.Header.Get("X-Chat-Registration-Token") == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//get files from multipart form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm
	if files == nil {
		http.Error(w, "No files found", http.StatusBadRequest)
		return
	}

	if len(files.File) == 0 {
		http.Error(w, "No files found", http.StatusBadRequest)
		return
	}

	var chatinfo chatInfo

	chatinfo.chatID, err = f.returnChatID(r.Header.Get("X-Chat-Registration-Token"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	f.telegramapiInst.SendFileToChat(files.File, chatinfo.chatID)

	// Save file to disk
	// f.saveFile(file)

}

func (f *FilesReceriverClient) returnChatID(token string) (int64, error) {
	f.sncmtx.Lock()
	defer f.sncmtx.Unlock()
	//check if chat is in cache
	val, ok := f.tokenCache[token]
	if !ok {
		//get chatid from token
		chatinfo, err := f.telegramapiInst.db.FindChatIDByToken(token)
		if err != nil {
			return 0, err
		}
		f.tokenCache[token] = chatinfo.chatID
		val = chatinfo.chatID
	}

	if val == 0 {
		return 0, fmt.Errorf("chat not found")
	}

	return val, nil
}
