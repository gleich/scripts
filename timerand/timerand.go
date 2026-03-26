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
	if err != nil || count <= 0 {
		timber.Fatal(err, "invalid count")
	}

	fmt.Print("Total time (minutes): ")
	scanner.Scan()
	totalMinutes, err := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)
	if err != nil || totalMinutes <= 0 {
		timber.Fatal(err, "invalid total time")
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
		start := time.Now()
		timber.Info(fmt.Sprintf("[%d/%d]", i+1, count), timber.A("waiting_for", util.FormatDuration(waitTime)))

		time.Sleep(waitTime)
		timber.DoneSince(start, "sending notification")

		err = exec.Command(
			"osascript", "-e",
			fmt.Sprintf(
				`display notification %q with title %q`,
				fmt.Sprintf("%d of %d", i+1, count),
				"timerand",
			)).Run()
		if err != nil {
			timber.Fatal(err, "failed to display notification")
		}

		for range 5 {
			err = exec.Command("afplay", "/System/Library/Sounds/Ping.aiff").Run()
			if err != nil {
				timber.Fatal(err, "failed to play ping sound")
			}
		}
	}

	timber.Done("all done")
}
