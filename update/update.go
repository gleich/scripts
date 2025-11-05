package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"go.mattglei.ch/timber"
)

type command struct {
	binary string
	args   []string
}

func (c command) String() string {
	s := c.binary
	if len(c.args) != 0 {
		s = fmt.Sprintf("%s %s", s, strings.Join(c.args, " "))
	}
	return s
}

var commands = []command{
	{binary: "brew", args: []string{"update"}},
	{binary: "brew", args: []string{"upgrade"}},
	{binary: "brew", args: []string{"cleanup", "-s"}},
	{binary: "rustup", args: []string{"update"}},
	{binary: "cargo", args: []string{"install-update", "-a"}},
	{binary: "code", args: []string{"--update-extensions"}},
}

func main() {
	timber.Timezone(time.Local)
	timber.TimeFormat("03:04:05")

	if slices.Contains(os.Args, "--dots") {
		commands = append(commands, command{binary: "dots"})
	}

	home, err := os.UserHomeDir()
	if err != nil {
		timber.Fatal(err, "failed to get user's home directory")
	}
	filepath := filepath.Join(home, ".update", "time.txt")
	now := time.Now()
	checkTime(now, filepath)
	writeTime(now, filepath)

	start := time.Now()
	elapsedTimes := []time.Duration{}
	for _, cmd := range commands {
		execStart := time.Now()

		timber.Info("running", cmd)

		cmdExec := exec.Command(cmd.binary, cmd.args...)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		cmdExec.Stdin = os.Stdin

		err := cmdExec.Run()
		if err != nil {
			timber.Fatal(err, "failed to run command")
		}

		elapsed := time.Since(execStart)
		elapsedTimes = append(elapsedTimes, elapsed)
		timber.Done("running", cmd, "in", elapsed)

	}

	fmt.Println()
	timber.Done("executed", len(commands), "commands in", time.Since(start))
	timber.Done("breakdown:")
	for i, cmd := range commands {
		fmt.Printf("\t%s: %s\n", cmd, elapsedTimes[i])
	}
}
