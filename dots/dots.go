package main

import (
	"fmt"
	"os"
	"os/exec"
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

	fmt.Println()
	timber.Info("copying folders")
	for parentDir, dirs := range folders {
		for _, dir := range dirs {
			syspath, dotspath, err := paths(parentDir, dir)
			if err != nil {
				timber.Fatal(err, "failed to get paths for", parentDir, dir)
			}
			err = os.CopyFS(dotspath, os.DirFS(syspath))
			if err != nil {
				timber.Fatal(err, "failed to copy", syspath)
			}
			timber.Done("copied", filepath.Join(parentDir, dir))
		}
	}

	fmt.Println()
	timber.Info("running commands")
	for _, command := range commands {
		out, err := exec.Command(command.cmd[0], command.cmd[1:]...).Output()
		if err != nil {
			timber.Fatal(err, "failed to run", command.cmd)
		}
		dotspath := filepath.Join(REPO_DIR, command.filename)
		err = os.WriteFile(dotspath, out, 0644)
		if err != nil {
			timber.Fatal(err, "failed to write output of command to", dotspath)
		}
		timber.Done("ran", command.name, "command")
	}

	fmt.Println()

	out, err := exec.Command("neofetch", "--stdout").Output()
	if err != nil {
		timber.Fatal(err, "failed to run neofetch command")
	}
	err = os.WriteFile(
		filepath.Join(REPO_DIR, "README.md"),
		fmt.Appendf([]byte{}, "# dots\n```txt\n%s\n```", out),
		0644,
	)
	if err != nil {
		timber.Fatal(err, "failed to write changes to README")
	}
	timber.Done("wrote neofetch summary to readme")

	fmt.Println()
	timber.Done("uploading changes")
	cmds := map[string]*exec.Cmd{
		"staged changes":    exec.Command("git", "add", "."),
		"committed changes": exec.Command("git", "commit", "-m", "chore: update"),
		"pushed changes":    exec.Command("git", "push"),
	}

	for msg, cmd := range cmds {
		cmd.Dir = REPO_DIR
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			if msg == "committed changes" {
				if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
					timber.Done("no changes to commit")
					return
				}
			}
			timber.Error(err, "failed to run", cmd.Args)
		}
		timber.Done(msg)
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
