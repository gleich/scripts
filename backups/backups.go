package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mattglei.ch/timber"
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
		name:     "Caprover",
		prefix:   "caprover-backup",
		suffix:   ".tar",
		filename: "caprover.tar",
	},
	{
		name:     "Strava",
		prefix:   "export_",
		suffix:   ".zip",
		length:   19,
		filename: "strava.zip",
	},
	{
		name:     "GitHub",
		suffix:   ".tar.gz",
		length:   43,
		filename: "github.tar.gz",
	},
	{
		name:     "Goodnotes",
		prefix:   "Backup ",
		suffix:   ".zip",
		length:   21,
		filename: "goodnotes.zip",
	},
	{
		name:     "Yamaha N800A",
		prefix:   "MC_backup_R-N800A",
		suffix:   ".dat",
		length:   21,
		filename: "yamaha-n800a.dat",
	},
	{
		name:     "Uniden R8",
		prefix:   "R8UserSetting",
		suffix:   ".bin",
		filename: "uniden_r8.bin",
	},
}

func main() {
	timber.Timezone(time.Local)
	timber.TimeFormat("03:04:05")

	home, err := os.UserHomeDir()
	if err != nil {
		timber.Fatal(err, "failed to get home directory")
	}

	downloadsPath := filepath.Join(home, "Downloads")
	entires, err := os.ReadDir(downloadsPath)
	if err != nil {
		timber.Fatal(err, "failed to read files from downloads folder")
	}

	backedUp := 0
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
				defer func() {
					err = sourceFile.Close()
					if err != nil {
						timber.Fatal(err, "failed to close source file")
					}
				}()

				destFile, err := os.Create(destination)
				if err != nil {
					timber.Fatal(err, "failed to create destination file")
				}
				defer func() {
					err = destFile.Close()
					if err != nil {
						timber.Fatal(err, "failed to close destination file")
					}
				}()

				_, err = io.Copy(destFile, sourceFile)
				if err != nil {
					timber.Fatal(err, "failed to copy backup file to destination")
				}

				err = os.Remove(sourcePath)
				if err != nil {
					timber.Fatal(err, "failed to remove source file")
				}

				timber.Done("Moved", backup.name)
				backedUp++
			}
		}
	}
	timber.Done("Backed up", backedUp, "items")
}
