package sgs_server

import (
	"fmt"
	"sgs/sgs_conf"
)

type SgsServer struct {
	HttpServer *SgsHttpServer
	UdpServer  *SgsUdpServer
}

func NewSgsServer(cfg *sgs_conf.SgsConfig) *SgsServer {
	s := new(SgsServer)
	s.HttpServer = NewSgsHttpServer(cfg.HttpConfig)
	if s.HttpServer == nil {
		fmt.Println("failed to new Http server")
		return nil
	}
	s.UdpServer = NewSgsUdpServer(cfg.UdpConfig)
	return s
}

func (s *SgsServer) Start() error {
	var err error
	fmt.Println("server=", s)
	go s.HttpServer.Start()

	go s.UdpServer.Start()

	return err
}
