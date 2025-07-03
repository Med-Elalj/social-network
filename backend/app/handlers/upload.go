package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	auth "social-network/app/Auth"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		auth.JsRespond(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		auth.JsRespond(w, "File upload error: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create uploads directory with explicit permissions
	if err := os.MkdirAll("uploads", 0o755); err != nil {
		auth.JsRespond(w, "Could not create directory", http.StatusInternalServerError)
		return
	}

	// Secure the filename and create path
	filename := filepath.Base(handler.Filename) // Prevent directory traversal
	filePath := filepath.Join("uploads", filename)

	dst, err := os.Create(filePath)
	if err != nil {
		auth.JsRespond(w, "Unable to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		auth.JsRespond(w, "Error copying file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return forward-slash path for web compatibility
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"path": "/uploads/" + filename, // Note the forward slash
	})
}

// TODO POST HANDLER
