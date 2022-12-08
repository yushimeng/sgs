package main

import (
	"github.com/yushimeng/sgs/log"
)

func main() {
	defer log.Close()

	for i := 0; i < 100; i++ {
		log.Debug("This is Debug msg---%d", i)
		log.Trace("This is Trace msg---%d", i)
		log.Info("This is Info msg---%d", i)
		log.Error("This is Error msg---%d", i)
		log.Fatal("This is Fatal msg---%d", i)
		log.Fatal("------------------------------%d", i)
	}
}
