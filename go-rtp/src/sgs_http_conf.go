package rtp

import "fmt"

type ConfigHttpServer struct {
	HttpPort int
}

func (cfg *ConfigHttpServer) SetDefaultConf() {
	cfg.HttpPort = 8080
}

func (cfg *ConfigHttpServer) HttpServerConfigCheck() (err error) {
	if cfg.HttpPort < 1 || cfg.HttpPort > 65535 {
		return fmt.Errorf("HttpPort[%d] should be in [1, 65535]", cfg.HttpPort)
	}
	return err
}
