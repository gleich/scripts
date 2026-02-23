package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/pkg/browser"
	"go.mattglei.ch/scripts/internal/logger"
	"go.mattglei.ch/timber"
)

type Inputs struct {
	TeamID      string
	KeyID       string
	KeyFilename string
}

func main() {
	logger.Setup()

	var (
		inputs Inputs
		keys   []huh.Option[string]
	)

	entries, err := os.ReadDir(".")
	if err != nil {
		timber.Fatal(err, "failed to read directory")
	}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && strings.HasSuffix(name, ".p8") {
			keys = append(keys, huh.NewOption(name, name))
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Team ID").
				Validate(huh.ValidateLength(10, 10)).
				Value(&inputs.TeamID),
			huh.NewInput().
				Title("Key ID").
				Validate(huh.ValidateLength(10, 10)).
				Value(&inputs.KeyID),
			huh.NewSelect[string]().
				Title("Key File").
				Options(keys...).
				Value(&inputs.KeyFilename),
		),
	).WithTheme(huh.ThemeBase())
	err = form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return
		}
		timber.Fatal(err, "failed to ask questions")
	}

	appToken, expiration := generateAppToken(inputs)

	timber.Done("app token:", appToken)
	timber.Done("expires:", expiration.Format("January 2 2006"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveUserAuth(appToken))

	addr := ":8000"
	url := "http://localhost" + addr
	fmt.Println()
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	timber.Info("starting server for user token:", url)
	go func() {
		err = http.ListenAndServe(addr, mux)
		if err != nil {
			timber.Fatal(err, "failed to start http server")
		}
	}()
	err = browser.OpenURL(url)
	if err != nil {
		timber.Fatal(err, "failed to open", url, "in browser")
	}
	waitGroup.Wait()
}
