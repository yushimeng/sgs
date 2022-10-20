package rtp

import (
	"fmt"

	gcfg "gopkg.in/gcfg.v1"
)

type SgsConfig struct {
	Basic      ConfigBasic
	HttpConfig ConfigHttpServer
	UdpConfig  ConfigRtpServer
}

func SgsConfigLoad(confFile, confRoot string) (SgsConfig, error) {
	var err error
	var cfg SgsConfig

	fmt.Printf("confFile:%s, confRoot:%s\n", confFile, confRoot)

	cfg.SetDefaultConf()

	// read config from file
	err = gcfg.ReadFileInto(&cfg, confFile)
	if err != nil {
		fmt.Printf("failed to Read config file, err:%s\n", err.Error())
		return cfg, err
	}

	// check params
	err = cfg.Check()
	return cfg, err
}

func (cfg *SgsConfig) SetDefaultConf() {
	cfg.Basic.SetDefaultConf()
	cfg.HttpConfig.SetDefaultConf()
	cfg.UdpConfig.SetDefaultConf()
}

func (cfg *SgsConfig) Check() (err error) {
	err = cfg.Basic.BasicConfigCheck()
	if err != nil {
		return err
	}

	err = cfg.HttpConfig.HttpServerConfigCheck()
	if err != nil {
		return err
	}

	err = cfg.UdpConfig.UdpServerConfigCheck()
	if err != nil {
		return err
	}

	return err
}
