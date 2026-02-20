package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os/exec"
	"time"

	"go.mattglei.ch/scripts/util"
	"go.mattglei.ch/timber"
)

func main() {
	timber.Timezone(time.Local)
	timber.TimeFormat("03:04:05")

	for {
		const min = 2 * time.Minute
		const span = 3 * time.Minute
		n, err := rand.Int(rand.Reader, big.NewInt(int64(span)))
		if err != nil {
			timber.Fatal(err, "failed to get random number")
		}
		waitTime := min + time.Duration(n.Int64())
		waitTimeFormatted := util.FormatDuration(waitTime)

		start := time.Now()
		timber.Info("Waiting for", waitTimeFormatted)

		time.Sleep(waitTime)
		timber.DoneSince(start, "sending notification")

		err = exec.Command(
			"osascript", "-e",
			fmt.Sprintf(
				`display notification %q with title %q`,
				fmt.Sprintf("%s is up", waitTimeFormatted),
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
}
