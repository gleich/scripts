package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"go.mattglei.ch/scripts/internal/logger"
	"go.mattglei.ch/scripts/internal/util"
	"go.mattglei.ch/timber"
)

func main() {
	logger.Setup()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Count: ")
	scanner.Scan()
	count, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		timber.Fatal(err, "invalid count")
	}
	if count < 0 {
		timber.Fatal(fmt.Errorf("count must be at least 0"), "invalid count")
	}
	if count == 0 {
		runRandomTimer()
	}

	fmt.Print("Total time (minutes): ")
	scanner.Scan()
	totalMinutes, err := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)
	if err != nil {
		timber.Fatal(err, "invalid total time")
	}
	if totalMinutes <= 0 {
		timber.Fatal(fmt.Errorf("total time must be greater than 0"), "invalid total time")
	}
	totalTime := time.Duration(totalMinutes * float64(time.Minute))

	// Generate random weights and normalize to sum to totalTime
	weights := make([]int64, count)
	var weightSum int64
	for i := range count {
		n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
		if err != nil {
			timber.Fatal(err, "failed to get random number")
		}
		// avoid zero weights
		weights[i] = n.Int64() + 1
		weightSum += weights[i]
	}

	intervals := make([]time.Duration, count)
	var allocated time.Duration
	for i := range count {
		if i == count-1 {
			intervals[i] = totalTime - allocated
		} else {
			intervals[i] = time.Duration(int64(totalTime) * weights[i] / weightSum)
			allocated += intervals[i]
		}
	}

	for i, waitTime := range intervals {
		progress := fmt.Sprintf("[%d/%d]", i+1, count)
		notification := fmt.Sprintf("%d of %d", i+1, count)
		waitAndNotify(waitTime, progress, notification)
	}

	timber.Done("all done")
}

func runRandomTimer() {
	const minWait = time.Minute
	const waitRange = 3 * time.Minute

	for {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(waitRange)))
		if err != nil {
			timber.Fatal(err, "failed to get random number")
		}
		waitTime := minWait + time.Duration(n.Int64())
		waitTimeFormatted := util.FormatDuration(waitTime)
		waitAndNotify(waitTime, "", fmt.Sprintf("%s is up", waitTimeFormatted))
	}
}

func waitAndNotify(waitTime time.Duration, progress string, notification string) {
	start := time.Now()
	waitTimeFormatted := util.FormatDuration(waitTime)
	if progress == "" {
		timber.Infof("waiting for %s", waitTimeFormatted)
	} else {
		timber.Infof("%s waiting for %s", progress, waitTimeFormatted)
	}

	time.Sleep(waitTime)
	timber.DoneSince(start, "sending notification")

	err := exec.Command(
		"osascript", "-e",
		fmt.Sprintf(
			`display notification %q with title %q`,
			notification,
			"timerand",
		)).Run()
	if err != nil {
		timber.Fatal(err, "failed to display notification")
	}

	for range 2 {
		err = exec.Command("afplay", "/System/Library/Sounds/Ping.aiff").Run()
		if err != nil {
			timber.Fatal(err, "failed to play ping sound")
		}
	}
}
