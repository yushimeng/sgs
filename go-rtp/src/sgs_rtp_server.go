package rtp

import (
	"net"
	"sgs/common"
)

type SgsRtpServer struct {
	server       *SgsServerManager
	Config       ConfigRtpServer
	udpTransport *common.SgsUdpTransport
	muxerMap     map[uint32]*RtmpMuxer
	// map ssrc -> rtp_stream
}

func NewSgsRtpServer(svr *SgsServerManager, cfg ConfigRtpServer) *SgsRtpServer {
	rtpServer := new(SgsRtpServer)
	rtpServer.server = svr
	rtpServer.Config = cfg
	rtpServer.udpTransport = common.NewSgsUdpTransport(cfg.UdpPort)
	rtpServer.udpTransport.SetHandle(rtpServer)
	return rtpServer
}

func (rtpServer *SgsRtpServer) OnRecv(data []byte, len int, from *net.UDPAddr) {
	pkt := NewRtpPacket()
	err := pkt.Parse(data, len, from)
	if err != nil {
		pkt.FreePacket()
		return
	}

	muxer, err := rtpServer.FetchOrCreateMuxer(pkt.Ssrc())
	if err != nil {
		pkt.FreePacket()
		return
	}

	muxer.Jitter.InsertPacket(pkt)
	return
}

func (rtpServer *SgsRtpServer) FetchOrCreateMuxer(ssrc uint32) (muxer *RtmpMuxer, err error) {
	if muxer, ok := muxerMap[ssrc]; ok {
		return muxer, ok
	}

	muxer = new(RtmpMuxer)
	return muxer, err
}
