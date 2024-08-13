package bot

import "net/http"

type FilesReceriverClient struct {
}

func (f *FilesReceriverClient) RecieveFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save file to disk
		// f.saveFile(file)
	}

}
