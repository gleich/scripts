package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"go.mattglei.ch/timber"
)

//go:embed auth.html
var authHTML string

func serveUserAuth(appToken string) http.HandlerFunc {
	page := strings.ReplaceAll(authHTML, "AMTOKENS_APP_TOKEN", fmt.Sprintf("'%s'", appToken))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("Cache-Control", "no-store")

		_, err := w.Write([]byte(page))
		if err != nil {
			timber.Fatal(err, "failed to write html to response")
		}
	}
}
