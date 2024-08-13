package bot

import "net/http"

type FilesReceriverClient struct {
	telegramapiInst *TelegramUpdates
}

func NewFilesReceriverClient(t *TelegramUpdates) *FilesReceriverClient {
	return &FilesReceriverClient{
		telegramapiInst: t,
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

	f.telegramapiInst.SendFileToChat(files.File, r.Header.Get("X-Chat-Registration-Token"))

	// Save file to disk
	// f.saveFile(file)

}
