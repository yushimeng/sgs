package rtp

import "fmt"

type ConfigRtpServer struct {
	UdpPort int
}

func (cfg *ConfigRtpServer) SetDefaultConf() {
	cfg.UdpPort = 8060
}

func (cfg *ConfigRtpServer) UdpServerConfigCheck() (err error) {
	if cfg.UdpPort < 1 || cfg.UdpPort > 65535 {
		return fmt.Errorf("UdpPort[%d] should be in [1, 65535]", cfg.UdpPort)
	}

	return err
}
