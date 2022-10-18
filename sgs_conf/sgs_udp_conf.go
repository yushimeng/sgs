package sgs_conf

import "fmt"

type ConfigUdpServer struct {
	UdpPort int
}

func (cfg *ConfigUdpServer) SetDefaultConf() {
	cfg.UdpPort = 8060
}
func (cfg *ConfigUdpServer) UdpServerConfigCheck() (err error) {
	if cfg.UdpPort < 1 || cfg.UdpPort > 65535 {
		return fmt.Errorf("UdpPort[%d] should be in [1, 65535]", cfg.UdpPort)
	}

	return err
}
