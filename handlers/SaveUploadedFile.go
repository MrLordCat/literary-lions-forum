package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func SaveUploadedFile(r *http.Request, formField string) (string, error) {
	file, header, err := r.FormFile(formField)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil 
		}
		return "", err
	}
	defer file.Close()

	if header == nil || header.Filename == "" {
		return "", nil 
	}

	filePath := filepath.Join("uploads", header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
