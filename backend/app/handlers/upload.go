package handlers

import (
	"bytes"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"social-network/app/logs"
	"social-network/app/modules"
)

const (
	maxUploadSize = 4 * (1 << 20) // 4 MiB
	uploadDir     = "../front-end/public/uploads"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "No file found in form"}`)
		return
	}
	defer file.Close()

	// Read into memory buffer for analysis
	tmpData, err := io.ReadAll(io.LimitReader(file, maxUploadSize+1))
	if err != nil {
		http.Error(w, `{"error":"failed to read file"}`, http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Failed to read file"}`)
		return
	}

	if len(tmpData) > maxUploadSize {
		http.Error(w, `{"error":"file too large"}`, http.StatusBadRequest)
		return
	}

	sniff := tmpData
	if len(sniff) > 512 {
		sniff = sniff[:512]
	}
	mediaType := http.DetectContentType(sniff)

	reader := bytes.NewReader(tmpData)

	if !modules.IsValidImage(reader) {
		http.Error(w, `{"error":"invalid image format"}`, http.StatusBadRequest)
		return
	}

	//make a unique filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		exts, _ := mime.ExtensionsByType(mediaType)
		if len(exts) > 0 {
			ext = exts[0]
		}
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	//save the file
	if err = os.WriteFile(filepath.Join(uploadDir, filename), tmpData, 0o644); err != nil {
		logs.ErrorLog.Printf("failed to save file: %v", err)
		http.Error(w, `{"error":"failed to save file"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"path": "/uploads/%s"}`, filename)
}
