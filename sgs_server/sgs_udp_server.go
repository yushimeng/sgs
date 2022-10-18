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

func (this *SgsUdpServer) Start() (err error) {
	address := fmt.Sprintf(":%d", this.Config.UdpPort)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	if err != nil {
		fmt.Println("read from connect failed, err:" + err.Error())
		util.AbnormalExit()
	}

	for {
		data := make([]byte, 1500)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err == nil {
			go udpProcess(data, n, remoteAddr, conn)
			// conn.WriteToUDP(data[:n], remoteAddr)
		} else {
			log.Logger.Error("failed to read from udp port,%s", err.Error())
			break
		}
	}

	return err
}

// UDP goroutine 实现并发读取UDP数据
func udpProcess(data []byte, len int, remoteAddr *net.UDPAddr, conn *net.UDPConn) {

	fmt.Println(data)
	conn.WriteToUDP(data[:len], remoteAddr)
}
