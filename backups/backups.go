package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
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

var grey = lipgloss.NewStyle().Foreground(lipgloss.Color("#4e4e4e"))

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

	count := 0
	for _, backup := range backups {
		backedUp := false
		for _, entry := range entires {
			name := entry.Name()
			if !entry.IsDir() && strings.HasPrefix(name, backup.prefix) &&
				strings.HasSuffix(
					name,
					backup.suffix,
				) && (backup.length == 0 || backup.length == len(name)) {

				var (
					destination = filepath.Join(
						home,
						"Library/Mobile Documents/com~apple~CloudDocs/Important/exports",
						backup.filename,
					)
					source = filepath.Join(downloadsPath, name)
				)

				err = os.Rename(source, destination)
				if err != nil {
					timber.Fatal(err, "failed to move", source, "to", destination)
				}

				timber.Done("Moved", backup.name)
				count++
				backedUp = true
				break
			}
		}
		if !backedUp {
			timber.Info(grey.Render(fmt.Sprintf("%s not backed up", backup.name)))
		}
	}
	fmt.Println()
	timber.Done("Backed up", count, "items")
}
