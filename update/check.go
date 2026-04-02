package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"go.mattglei.ch/timber"
)

const timeFormat = "2006-01-02"

func checkTime(now time.Time, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	bin, err := os.ReadFile(path)
	if err != nil {
		timber.Fatal(err, "failed to read from", timber.A("path", path))
	}
	t, err := time.ParseInLocation(timeFormat, strings.TrimSpace(string(bin)), time.Local)
	if err != nil {
		timber.Fatal(err, "failed to parse time")
	}
	t = t.Local()
	today := t.Format(timeFormat) == now.Format(timeFormat)
	if !today {
		return
	}

	var runAgain bool
	err = huh.NewConfirm().
		Title("Update has already been ran today. Run again?").
		Value(&runAgain).
		WithTheme(huh.ThemeBase()).
		Run()
	if err != nil {
		timber.Fatal(err, "failed to get user's input")
	}
	if !runAgain {
		os.Exit(0)
	}
}

func writeTime(now time.Time, path string) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		timber.Fatal(err, "failed to make directory", timber.A("dir", dir))
	}
	err = os.WriteFile(path, []byte(now.Format(timeFormat)), 0600)
	if err != nil {
		timber.Fatal(err, "failed to write to file")
	}
}
