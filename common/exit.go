package common

import (
	"os"

	"github.com/baidu/go-lib/log"
)

func AbnormalExit() {
	// waiting for logger finish jobs
	log.Logger.Close()
	// exit
	os.Exit(1)
}