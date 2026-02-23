package logger

import (
	"time"

	"go.mattglei.ch/timber"
)

func Setup() {
	timber.Timezone(time.Local)
	timber.TimeFormat("03:04:05")
}
