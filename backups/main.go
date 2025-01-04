package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gleich/lumber/v3"
)

type Backup struct {
	Prefix   *string `toml:"prefix"`
	Suffix   *string `toml:"suffix"`
	Filename *string `toml:"filename"`
}

func main() {
	setupLogger()
	lumber.Done("booted")

	home, err := os.UserHomeDir()
	if err != nil {
		lumber.Fatal(err, "failed to get home directory")
	}

	conf := readConfig(home)

	downloadsPath := filepath.Join(home, "Downloads")
	entires, err := os.ReadDir(downloadsPath)
	if err != nil {
		lumber.Fatal(err, "failed to read files from downloads folder")
	}
	for backupName, backup := range conf {
		for _, entry := range entires {
			name := entry.Name()
			if !entry.IsDir() && strings.HasPrefix(name, *backup.Prefix) &&
				strings.HasSuffix(name, *backup.Suffix) {
				destination := filepath.Join(
					home,
					"Library/Mobile Documents/com~apple~CloudDocs/Important/exports",
					*backup.Filename,
				)
				if _, err := os.Stat(destination); !errors.Is(err, os.ErrNotExist) {
					err = os.Remove(destination)
					if err != nil {
						lumber.Fatal(err, "failed to delete destination file")
					}
				}

				sourcePath := filepath.Join(downloadsPath, name)
				sourceFile, err := os.Open(sourcePath)
				if err != nil {
					lumber.Fatal(err, "failed to open source file")
				}
				defer sourceFile.Close()

				destFile, err := os.Create(destination)
				if err != nil {
					lumber.Fatal(err, "failed to create destination file")
				}
				defer destFile.Close()

				_, err = io.Copy(destFile, sourceFile)
				if err != nil {
					lumber.Fatal(err, "failed to copy backup file to destination")
				}

				err = os.Remove(sourcePath)
				if err != nil {
					lumber.Fatal(err, "failed to remove source file")
				}

				lumber.Done("Moved", backupName)
			}
		}
	}
}

func setupLogger() {
	nytime, err := time.LoadLocation("America/New_York")
	if err != nil {
		lumber.Fatal(err, "failed to load new york timezone")
	}
	lumber.SetTimezone(nytime)
	lumber.SetTimeFormat("01/02 03:04:05 PM MST")
}

func readConfig(home string) map[string]Backup {
	var backups map[string]Backup
	_, err := toml.DecodeFile(
		filepath.Join(home, "src/gleich/scripts/backups/config.toml"),
		&backups,
	)
	if err != nil {
		lumber.Error(err, "failed to parse config file")
	}
	return backups
}
