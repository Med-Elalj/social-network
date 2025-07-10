package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"social-network/app/logs"
	"social-network/app/modules"
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
	tmpData, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Failed to read file"}`)
		return
	}

	// Check image validity
	reader := strings.NewReader(string(tmpData))
	if !modules.IsValidImage(reader, header.Header.Get("Content-Type")) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid image format"}`)
		return
	}

	// Create a unique filename
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", timestamp, ext)

	// Save the file to disk
	err = os.WriteFile(filepath.Join("../../../front-end/public/uploads", filename), tmpData, 0o644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logs.ErrorLog.Println("Failed to save file:", err)
		fmt.Fprintf(w, `{"error": "Failed to save file"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"filename": "%s"}`, filename)
}
