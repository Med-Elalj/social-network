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

	auth "social-network/app/Auth"
	"social-network/app/logs"
	"social-network/app/modules"
)

const (
	maxUploadSize = 4 * (1 << 20) // 4 MiB
	uploadDir     = "../front-end/public/uploads"
)

type UploadResponse struct {
	Message string `json:"message"`
	Path    string `json:"path"`
	Code    int    `json:"code"`
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		auth.JsResponse(w, "No file found in form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read into memory buffer for analysis
	tmpData, err := io.ReadAll(io.LimitReader(file, maxUploadSize+1))
	if err != nil {
		auth.JsResponse(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	if len(tmpData) > maxUploadSize {
		auth.JsResponse(w, "File exceeds 4MB size limit", http.StatusRequestEntityTooLarge)
		return
	}

	sniff := tmpData
	if len(sniff) > 512 {
		sniff = sniff[:512]
	}
	mediaType := http.DetectContentType(sniff)

	reader := bytes.NewReader(tmpData)

	if !modules.IsValidImage(reader) {
		auth.JsResponse(w, "Invalid image format", http.StatusBadRequest)
		return
	}

	// make a unique filename
	now := time.Now()
	dateDir := now.Format("2006/01/02")
	fullUploadPath := filepath.Join(uploadDir, dateDir)

	if err := os.MkdirAll(fullUploadPath, 0o755); err != nil {
		logs.ErrorLog.Printf("failed to create upload directory: %v", err)
		auth.JsResponse(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		exts, _ := mime.ExtensionsByType(mediaType)
		if len(exts) > 0 {
			ext = exts[0]
		}
	}
	filename := fmt.Sprintf("%d%s", now.UnixNano(), ext)

	// save the file
	if err = os.WriteFile(filepath.Join(fullUploadPath, filename), tmpData, 0o644); err != nil {
		logs.ErrorLog.Printf("failed to save file: %v", err)
		auth.JsResponse(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	resp := UploadResponse{
		Message: "File uploaded successfully",
		Path:    fmt.Sprintf("/uploads/%s/%s", dateDir, filename),
		Code:    http.StatusOK,
	}

	auth.JsMapResponse(w, UploadResponseToMap(resp), http.StatusOK)
}

func UploadResponseToMap(resp UploadResponse) map[string]any {
	return map[string]any{
		"message": resp.Message,
		"path":    resp.Path,
		"code":    resp.Code,
	}
}
