package bot

import (
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

	f.sncmtx.Lock()
	//check if chat is in cache
	_, ok := f.tokenCache[r.Header.Get("X-Chat-Registration-Token")]
	if !ok {
		//get chatid from token
		chatinfo, err := f.telegramapiInst.db.FindChatIDByToken(r.Header.Get("X-Chat-Registration-Token"))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		f.tokenCache[r.Header.Get("X-Chat-Registration-Token")] = chatinfo.chatID
	}

	f.sncmtx.Unlock()

	f.telegramapiInst.SendFileToChat(files.File, chatinfo.chatID)

	// Save file to disk
	// f.saveFile(file)

}
