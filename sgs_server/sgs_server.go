package sgs_server

import (
	"fmt"
	"sgs/sgs_conf"

	"github.com/baidu/go-lib/log"
)

type SgsServer struct {
	Config     *sgs_conf.SgsConfig
	HttpServer *SgsHttpServer
	UdpServer  *SgsUdpServer
}

func NewSgsServer(cfg *sgs_conf.SgsConfig) *SgsServer {
	server := new(SgsServer)
	server.Config = cfg
	server.HttpServer = NewSgsHttpServer(server, cfg.HttpConfig)
	if server.HttpServer == nil {
		fmt.Println("failed to new Http server")
		return nil
	}

	server.UdpServer = NewSgsUdpServer(cfg.UdpConfig)
	if server.UdpServer == nil {
		log.Logger.Error("failed to new udp server")
	}
	return server
}

func (server *SgsServer) Start() {
	fmt.Println("server=", server)

	go server.HttpServer.Start()

	go server.UdpServer.Start()
}
