package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"pkg.mattglei.ch/timber"
)

type command struct {
	Binary    string
	Args      []string
	Directory string
}

var commands = []command{
	{Binary: "brew", Args: []string{"update"}},
	{Binary: "brew", Args: []string{"upgrade"}},
	{Binary: "brew", Args: []string{"cleanup", "-s"}},
	{Binary: "rustup", Args: []string{"update"}},
	{Binary: "fetch", Directory: "/Users/matt/src/gleich/dots"},
}

func main() {
	timber.SetTimezone(time.Local)
	timber.SetTimeFormat("03:04:05")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	start := time.Now()
	elapsedTimes := []time.Duration{}
	for _, cmd := range commands {
		execStart := time.Now()

		args := strings.Join(cmd.Args, " ")
		timber.Info("running", cmd.Binary, args)

		cmdExec := exec.CommandContext(ctx, cmd.Binary, cmd.Args...)
		if cmd.Directory != "" {
			cmdExec.Dir = cmd.Directory
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
		timber.Done("finished running", cmd.Binary, args, "in", elapsed)

	}

	fmt.Println()
	timber.Done("executed", len(commands), "commands in", time.Since(start))
	timber.Done("breakdown:")
	for i, cmd := range commands {
		fmt.Printf(
			"\t%s",
			cmd.Binary,
		)
		if len(cmd.Args) != 0 {
			fmt.Printf(" %s", strings.Join(cmd.Args, " "))
		}
		fmt.Printf(": %s\n", elapsedTimes[i])
	}
}
