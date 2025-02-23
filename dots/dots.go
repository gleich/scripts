package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mattglei.ch/timber"
)

const REPO_DIR = "/Users/matt/src/gleich/dots/"

type command struct {
	name     string
	cmd      []string
	filename string
}

var (
	DOTS_ROOT_DIR = filepath.Join(REPO_DIR, "matt")

	folders = map[string][]string{
		"~": {
			"Library/Application Support/Code/User/Snippets",
		},
		"~/.config": {
			"fish",
			"omf",
			"resin",
			"neofetch",
			"zathura",
			"raycast/script-commands",
			"nvim",
			"gh",
			"kitty",
		},
	}
	files = map[string][]string{
		"~": {
			".gitconfig",
			".hushlogin",
			".docker/config.json",
			"Library/Application Support/Code/User/keybindings.json",
			"Library/Application Support/Code/User/settings.json",
			".gnupg/gpg-agent.conf",
			".rustup/settings.toml",
			".cargo/.crates.toml",
		},
		"~/.config": {
			"gh/config.yml",
		},
	}
	commands = []command{
		{
			name:     "homebrew",
			cmd:      []string{"brew", "bundle", "dump", "--describe", "--file=-"},
			filename: "Brewfile",
		},
		{name: "vscode", cmd: []string{"code", "--list-extensions"}, filename: "extensions.txt"},
	}
)

func main() {
	timber.Timezone(time.Local)
	timber.TimeFormat("03:04:05")

	err := os.RemoveAll(DOTS_ROOT_DIR)
	if err != nil {
		timber.Fatal(err, "failed to reset root dots directory:", DOTS_ROOT_DIR)
	}

	timber.Info("copying files")
	for dir, filenames := range files {
		for _, filename := range filenames {
			syspath, dotspath, err := paths(dir, filename)
			if err != nil {
				timber.Fatal(err, "failed to get path for", dir, filename)
			}
			err = os.MkdirAll(filepath.Dir(dotspath), os.ModePerm)
			if err != nil {
				timber.Fatal(err, "failed to make parent dir for", dotspath)
			}
			data, err := os.ReadFile(syspath)
			if err != nil {
				timber.Fatal(err, "failed to read", syspath)
			}
			err = os.WriteFile(dotspath, data, 0644)
			if err != nil {
				timber.Fatal(err, "failed to write data to", dotspath)
			}
			timber.Done("copied", filepath.Join(dir, filename))
		}
	}
}

// get the system path and the dots repo path
func paths(dir string, fpath string) (string, string, error) {
	var sysdir, dotsdir string
	if strings.HasPrefix(dir, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", "", fmt.Errorf("%w: failed to get user's home dir", err)
		}
		cleanedDir := dir[1:]
		sysdir = filepath.Join(home, cleanedDir)
		dotsdir = filepath.Join(DOTS_ROOT_DIR, cleanedDir)
	} else {
		sysdir = dir
		dotsdir = filepath.Join(DOTS_ROOT_DIR, dir)
	}
	return filepath.Join(sysdir, fpath), filepath.Join(dotsdir, fpath), nil
}
