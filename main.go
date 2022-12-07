package main

import (
	"time"

	"github.com/sgs/log"
)

func main() {

	log.SetOutput("./sgs.log")
	log.Info("info msgs")
	ch := make(chan int)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			ch <- 1
		}
	}()
	<-ch
	log.Info("endmsg")
	glog.Close()
}
