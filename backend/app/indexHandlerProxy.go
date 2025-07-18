//go:build useproxy

package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func IndexHandler() http.Handler {
	target, err := url.Parse("http://localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	return httputil.NewSingleHostReverseProxy(target)
}
