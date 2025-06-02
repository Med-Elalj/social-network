//go:build !dev
// +build !dev

package sn

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

var (
	rootDir = ".front-end/dist"
	fs      = http.FileServer(http.Dir(rootDir))
)

func IndexHandler() http.Handler {
	fmt.Println("Using static file server at", rootDir)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean("/" + r.URL.Path) // ensure it starts with '/' for path.Clean

		// Disallow path traversal
		if strings.Contains(cleanPath, "..") {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}

		// Compute full path and ensure it stays within rootDir
		absPath, err := filepath.Abs(filepath.Join(rootDir, cleanPath))
		rootAbs, _ := filepath.Abs(rootDir)
		if err != nil || !strings.HasPrefix(absPath, rootAbs) {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}

		// Serve the sanitized request
		r.URL.Path = cleanPath
		fs.ServeHTTP(w, r)
	})
}
