package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	auth "social-network/app/Auth"
	db "social-network/app/modules"
)

const (
	maxFileSize   = 1 << 20 // 1 MB
	uploadDir     = "./server/data/images"
	formFieldName = "image"
)

// func EnsureUploadDir() error {
// 	return os.MkdirAll(uploadDir, os.ModePerm)
// }

func isValidImage(file io.ReadSeeker, contentType string) bool {
	var err error
	if strings.Contains(contentType, "svg") {
		err = validateSecureSVG(file)
	} else {
		_, _, err = image.Decode(file)
	}
	_, _ = file.Seek(0, 0)
	return err == nil
}

func UploadHandler(w http.ResponseWriter, r *http.Request, uid int) {
	// Limit request body size to maxFileSize
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize+1024)

	if err := r.ParseMultipartForm(maxFileSize + 1024); err != nil {
		auth.JsRespond(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile(formFieldName)
	if err != nil {
		auth.JsRespond(w, "No file found in form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read into memory buffer for analysis
	tmpData, err := io.ReadAll(io.LimitReader(file, maxFileSize+1))
	if err != nil {
		auth.JsRespond(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	if len(tmpData) > maxFileSize {
		auth.JsRespond(w, "File exceeds 1MB size limit", http.StatusRequestEntityTooLarge)
		return
	}

	// Check image validity
	reader := strings.NewReader(string(tmpData))
	if !isValidImage(reader, header.Header.Get("Content-Type")) {
		auth.JsRespond(w, "Invalid or unsupported image format", http.StatusBadRequest)
		return
	}

	// Compute hash for deduplication
	hash := sha256.Sum256(tmpData)
	hashStr := hex.EncodeToString(hash[:])

	// Use hash as filename to avoid duplication
	ext := filepath.Ext(header.Filename)
	filename := hashStr + ext
	destPath := filepath.Join(uploadDir, filename)

	if _, err := os.Stat(destPath); err == nil {
		existingData, err := os.ReadFile(destPath)
		if err == nil && string(existingData) == string(tmpData) {
			fmt.Fprintf(w, "Image uploaded successfully as %s\n", filename)
			return
		}
		// Find a unique filename with _%d suffix
		base := hashStr
		for i := 1; ; i++ {
			newFilename := fmt.Sprintf("%s_%d%s", base, i, ext)
			newDestPath := filepath.Join(uploadDir, newFilename)
			if _, err := os.Stat(newDestPath); os.IsNotExist(err) {
				if err := os.WriteFile(newDestPath, tmpData, 0o644); err != nil {
					auth.JsRespond(w, "Failed to save file", http.StatusInternalServerError)
					return
				}
				filename = newFilename
				break
			}
			// If file exists, check if it's the same
			existingData, err := os.ReadFile(newDestPath)
			if err == nil && string(existingData) == string(tmpData) {
				break
			}
		}
	} else {
		if err := os.WriteFile(destPath, tmpData, 0o644); err != nil {
			auth.JsRespond(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
	}
	// Add file info to user_files in db: uid, filename, size in bytes
	if err := db.AddUserFile(uid, filename, len(tmpData)); err != nil {
		// Attempt to delete the uploaded file if DB update fails
		err = os.Remove(destPath)
		if err != nil {
			log.Printf("Failed to delete file %s after DB update failure: %v", destPath, err)
		}
		auth.JsRespond(w, "Failed to update user files in database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Image uploaded successfully as %s\n", filename)
}
