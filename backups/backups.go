package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pkg.mattglei.ch/timber"
)

type backup struct {
	name     string
	prefix   string
	suffix   string
	length   int // optional
	filename string
}

var backups = []backup{
	{
		name:     "caprover",
		prefix:   "caprover-backup",
		suffix:   ".tar",
		filename: "caprover.tar",
	},
	{
		name:     "strava",
		prefix:   "export_",
		suffix:   ".zip",
		length:   19,
		filename: "strava.zip",
	},
	{
		name:     "github",
		suffix:   ".tar.gz",
		length:   43,
		filename: "github.tar.gz",
	},
	{
		name:     "goodnotes",
		prefix:   "Backup ",
		suffix:   ".zip",
		length:   21,
		filename: "goodnotes.zip",
	},
}

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	home, err := os.UserHomeDir()
	if err != nil {
		timber.Fatal(err, "failed to get home directory")
	}

	downloadsPath := filepath.Join(home, "Downloads")
	entires, err := os.ReadDir(downloadsPath)
	if err != nil {
		timber.Fatal(err, "failed to read files from downloads folder")
	}
	for _, backup := range backups {
		for _, entry := range entires {
			name := entry.Name()
			if !entry.IsDir() && strings.HasPrefix(name, backup.prefix) &&
				strings.HasSuffix(
					name,
					backup.suffix,
				) && (backup.length == 0 || backup.length == len(name)) {
				destination := filepath.Join(
					home,
					"Library/Mobile Documents/com~apple~CloudDocs/Important/exports",
					backup.filename,
				)
				if _, err := os.Stat(destination); !errors.Is(err, os.ErrNotExist) {
					err = os.Remove(destination)
					if err != nil {
						timber.Fatal(err, "failed to delete destination file")
					}
				}

				sourcePath := filepath.Join(downloadsPath, name)
				sourceFile, err := os.Open(sourcePath)
				if err != nil {
					timber.Fatal(err, "failed to open source file")
				}
				defer sourceFile.Close()

				destFile, err := os.Create(destination)
				if err != nil {
					timber.Fatal(err, "failed to create destination file")
				}
				defer destFile.Close()

				_, err = io.Copy(destFile, sourceFile)
				if err != nil {
					timber.Fatal(err, "failed to copy backup file to destination")
				}

				err = os.Remove(sourcePath)
				if err != nil {
					timber.Fatal(err, "failed to remove source file")
				}

				timber.Done("moved", backup.name)
			}
		}
	}
}
