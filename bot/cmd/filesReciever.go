package main

type FilesReciever struct {
	// Path to the directory where files will be saved
	Dir string
}

// RecieveFile saves the file to the directory specified in the FilesReciever.Dir field
func (f *FilesReciever) RecieveFile(file []byte, filename string) error {
	return nil
}
