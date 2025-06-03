//go:build dev
// +build dev

package sn

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func IndexHandler() http.Handler {
	// Assumes Next.js dev server is running at localhost:3000
	target, err := url.Parse("http://localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Using Next.js dev server at", target.String())
	return httputil.NewSingleHostReverseProxy(target)
}
