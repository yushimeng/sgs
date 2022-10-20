package rtp

import (
	"fmt"

	"github.com/baidu/go-lib/log"
)

type SgsServerManager struct {
	Config     *SgsConfig
	HttpServer *SgsHttpServer
	RtpServer  *SgsRtpServer
}

func NewSgsServer(cfg *SgsConfig) *SgsServerManager {
	serverManager := new(SgsServerManager)
	serverManager.Config = cfg
	serverManager.HttpServer = NewSgsHttpServer(serverManager, cfg.HttpConfig)
	if serverManager.HttpServer == nil {
		fmt.Println("failed to new Http server")
		return nil
	}

	serverManager.RtpServer = NewSgsRtpServer(serverManager, cfg.UdpConfig)
	if serverManager.RtpServer == nil {
		log.Logger.Error("failed to new udp server")
	}
	return serverManager
}

func (server *SgsServerManager) Start() {
	fmt.Println("server=", server)

	go server.HttpServer.Start()

	go server.RtpServer.Start()
}
