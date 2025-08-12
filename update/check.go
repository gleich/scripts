package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mattglei.ch/timber"
)

const timeFormat = "2006-01-02"

func checkTime(now time.Time, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	bin, err := os.ReadFile(path)
	if err != nil {
		timber.Fatal(err, "failed to read from", path)
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

	timber.Warning("Update has already been ran today. Do you want to run it again (y/n)?")
	var response string
	_, err = fmt.Scan(&response)
	if err != nil {
		timber.Fatal(err, "failed to get user's input")
	}
	if !strings.Contains(response, "y") {
		os.Exit(0)
	}
}

func writeTime(now time.Time, path string) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		timber.Fatal(err, "failed to make directory:", dir)
	}
	err = os.WriteFile(path, []byte(now.Format(timeFormat)), 0655)
	if err != nil {
		timber.Fatal(err, "failed to write to file")
	}
}
