package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"pkg.mattglei.ch/timber"
)

type Backup struct {
	Prefix   *string `toml:"prefix"`
	Suffix   *string `toml:"suffix"`
	Length   int     `toml:"length"` // optional
	Filename *string `toml:"filename"`
}

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	home, err := os.UserHomeDir()
	if err != nil {
		timber.Fatal(err, "failed to get home directory")
	}

	conf := readConfig(home)

	downloadsPath := filepath.Join(home, "Downloads")
	entires, err := os.ReadDir(downloadsPath)
	if err != nil {
		timber.Fatal(err, "failed to read files from downloads folder")
	}
	for backupName, backup := range conf {
		for _, entry := range entires {
			name := entry.Name()
			if !entry.IsDir() && strings.HasPrefix(name, *backup.Prefix) &&
				strings.HasSuffix(
					name,
					*backup.Suffix,
				) && (backup.Length == 0 || backup.Length == len(name)) {
				destination := filepath.Join(
					home,
					"Library/Mobile Documents/com~apple~CloudDocs/Important/exports",
					*backup.Filename,
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

				timber.Done("Moved", backupName)
			}
		}
	}
}

func readConfig(home string) map[string]Backup {
	var backups map[string]Backup
	_, err := toml.DecodeFile(
		filepath.Join(home, "src/gleich/scripts/backups/config.toml"),
		&backups,
	)
	if err != nil {
		timber.Error(err, "failed to parse config file")
	}
	return backups
}
