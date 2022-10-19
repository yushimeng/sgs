package main

import (
	"flag"
	"fmt"
	"path"
	"runtime"
	"sgs/sgs_conf"
	"sgs/sgs_server"
	"time"

	"github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/log/log4go"
)

var (
	help        = flag.Bool("h", false, "to show help")
	confRoot    = flag.String("c", "conf/", "config file rootPath")
	logRoot     = flag.String("l", "log/", "log file rootPath")
	showVersion = flag.Bool("v", false, "show version")
	stdOut      = flag.Bool("s", false, "to show log in stdout")
	debug       = flag.Bool("d", false, "debug mode")
)

var version string

func main() {
	var err error
	var logSwitch string
	var config sgs_conf.SgsConfig

	// initialize args
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	if *showVersion {
		fmt.Println(version)
		return
	}
	if *debug {
		logSwitch = "Debug"
	} else {
		logSwitch = "INFO"
	}

	// initialize log
	log4go.SetLogBufferLength(10000)
	log4go.SetLogWithBlocking(false)
	log4go.SetLogFormat(log4go.FORMAT_DEFAULT_WITH_PID)
	log4go.SetSrcLineForBinLog(false)
	fmt.Println("sgs log init")
	err = log.Init("sgs", logSwitch, *logRoot, *stdOut, "midnight", 7)
	if err != nil {
		fmt.Printf("sgs: err in log.Init():%s\n", err.Error())
		return
	}

	log.Logger.Info("sgs[version:%s] start", version)

	// initialize config
	confPath := path.Join(*confRoot, "sgs.conf")
	config, err = sgs_conf.SgsConfigLoad(confPath, *confRoot)
	if err != nil {
		fmt.Println("sgs conf load failed")
		log.Logger.Error("main(): in BfeConfigLoad():%s", err.Error())
		return
	}

	// maximum number of CPUs (GOMAXPROCS) defaults to runtime.CPUNUM
	// if running on machine, or CPU quota if running on container
	// (with the help of "go.uber.org/automaxprocs").
	// here, we change maximum number of cpus if the MaxCpus is positive.
	if config.Basic.MaxCpus > 0 {
		runtime.GOMAXPROCS(config.Basic.MaxCpus)
	}

	sgsServer := sgs_server.NewSgsServer(&config)
	sgsServer.Start()

	for {
		time.Sleep(time.Second)
	}
}
