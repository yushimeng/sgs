package common

import (
	"fmt"
	"net"

	"github.com/baidu/go-lib/log"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/common"
)

type HandleFunc interface {
	OnRecv(data []byte, len int, addr *net.UDPAddr)
}

type SgsUdpTransport struct {
	Handle HandleFunc
	port   int
}

func NewSgsUdpTransport(port int) (tr *SgsUdpTransport) {
	tr = new(SgsUdpTransport)
	tr.Handle = tr
	return tr
}

func (udpTransport *SgsUdpTransport) OnRecv(data []byte, len int, addr *net.UDPAddr) {
	fmt.Printf("TransportUDP: no registered upper layer RTP packet handler\n")
}

func (udpTransport *SgsUdpTransport) SetHandle(h HandleFunc) {
	udpTransport.Handle = h
}

func (udpTransport *SgsUdpTransport) Start() {
	address := fmt.Sprintf(":%d", udpTransport.port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		common.AbnormalExit()
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		common.AbnormalExit()
		return
	}
	defer conn.Close()

	for {
		data := make([]byte, 1500)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err == nil {
			if udpTransport.Handle != nil {
				go udpTransport.Handle.OnRecv(data, n, remoteAddr, conn)
			}
		} else {
			log.Logger.Error("failed to read from udp port,%s", err.Error())
			common.AbnormalExit()
			return
		}
	}
}
