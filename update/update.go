package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"go.mattglei.ch/timber"
)

type command struct {
	binary    string
	args      []string
	directory string
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
	{binary: "fetch", directory: "/Users/matt/src/gleich/dots"},
}

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	start := time.Now()
	elapsedTimes := []time.Duration{}
	for _, cmd := range commands {
		execStart := time.Now()

		timber.Info("running", cmd)

		cmdExec := exec.Command(cmd.binary, cmd.args...)
		if cmd.directory != "" {
			cmdExec.Dir = cmd.directory
		}
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		cmdExec.Stdin = os.Stdin

		err := cmdExec.Run()
		if err != nil {
			timber.Fatal(err, "failed to run command")
		}

		elapsed := time.Since(execStart)
		elapsedTimes = append(elapsedTimes, elapsed)
		timber.Done("finished running", cmd, "in", elapsed)

	}

	fmt.Println()
	timber.Done("executed", len(commands), "commands in", time.Since(start))
	timber.Done("breakdown:")
	for i, cmd := range commands {
		fmt.Printf("\t%s: %s\n", cmd, elapsedTimes[i])
	}
}
