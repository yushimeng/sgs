package sgs_server

import (
	"fmt"
	"net"
	"sgs/sgs_conf"
	"sgs/util"

	"github.com/baidu/go-lib/log"
)

type SgsUdpServer struct {
	Config sgs_conf.ConfigUdpServer
}

func NewSgsUdpServer(cfg sgs_conf.ConfigUdpServer) (s *SgsUdpServer) {
	s = new(SgsUdpServer)
	s.Config = cfg
	return s
}

func (udpServer *SgsUdpServer) Start() {
	address := fmt.Sprintf(":%d", udpServer.Config.UdpPort)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		util.AbnormalExit()
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		util.AbnormalExit()
		return
	}
	defer conn.Close()

	for {
		data := make([]byte, 1500)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err == nil {
			go udpProcess(data, n, remoteAddr, conn)
		} else {
			log.Logger.Error("failed to read from udp port,%s", err.Error())
			util.AbnormalExit()
			return
		}
	}
}

// UDP goroutine 实现并发读取UDP数据
func udpProcess(data []byte, len int, remoteAddr *net.UDPAddr, conn *net.UDPConn) {
	str := string(data[:len])
	fmt.Println(str)
	conn.WriteToUDP(data[:len], remoteAddr)
	fmt.Println("send back end...")
}
