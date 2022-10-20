package rtp

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
)

const (
	freeRtpPacketListNum = 10
	defaultMTULength     = 1500
	rtpHeaderLength      = 12
)

// RFC 3550 RTP Packet Format. refs: https://www.cl.cam.ac.uk/~jac22/books/mm/book/node159.html
//
//	0                   1                   2                   3
//	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |V=2|P|X|  CC   |M|     PT      |       sequence number         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                           timestamp                           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |           synchronization source (SSRC) identifier            |
// +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
// |            contributing source (CSRC) identifiers             |
// |                             ....                              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

// Marker returns the marker as bool in host order.
func (pkt *RtpPacket) Parse(data []byte, n int, addr *net.UDPAddr) error {
	var err error
	if n < 13 {
		return errors.New("Parse: rtp length less than 13")
	}
	copy(pkt.buf, data)
	pkt.inUse = n
	pkt.conn = addr
	pkt.fromAddr.IPAddr = addr.IP
	pkt.fromAddr.DataPort = addr.Port
	return err
}

// Marker returns the marker as bool in host order.
func (pkt *RtpPacket) Marker() bool {
	return (pkt.buf[markerOffset] & markerBit) == markerBit
}

// PayloadType return the payload type value from RTP packet header.
func (rp *RtpPacket) PayloadType() byte {
	return rp.buf[markerOffset] & ptMask
}

// Version returns the version as uint32 in host order.
func (pkt *RtpPacket) Version() uint8 {
	return (pkt.buf[versionOffset] & version2Bit)
}

// Version returns the version as uint32 in host order.
func (pkt *RtpPacket) Padding() bool {
	return (pkt.buf[paddingBitOffset] & paddingBit) == paddingBit
}

// Version returns the version as uint32 in host order.
func (pkt *RtpPacket) ExtensionBit() bool {
	return (pkt.buf[extensionBitOffset] & extensionBit) == extensionBit
}

// Version returns the version as uint32 in host order.
func (pkt *RtpPacket) CsrcCount() uint8 {
	return (pkt.buf[csrcCountOffset] & ccMask)
}

// Version returns the version as uint32 in host order.
func (pkt *RtpPacket) Sequence() uint16 {
	return binary.BigEndian.Uint16(pkt.buf[sequenceOffset:])
}

// Version returns the version as uint32 in host order.
func (pkt *RtpPacket) Timestamp() uint32 {
	return binary.BigEndian.Uint32(pkt.buf[timestampOffset:])
}

// Version returns the version as uint32 in host order.
func (pkt *RtpPacket) Ssrc() uint32 {
	return binary.BigEndian.Uint32(pkt.buf[ssrcOffsetRtp:])
}

// ExtensionLength returns the full length in bytes of RTP packet extension (including the main extension header).
func (pkt *RtpPacket) ExtensionLength() (length int) {
	if !pkt.ExtensionBit() {
		return 0
	}
	offset := int16(pkt.CsrcCount()*4 + rtpHeaderLength) // offset to extension header 32bit word
	offset += 2
	length = int(binary.BigEndian.Uint16(pkt.buf[offset:])) + 1 // +1 for the main extension header word
	length *= 4
	return
}

// CsrcList returns the list of CSRC values as uint32 slice in host horder
func (rp *RtpPacket) CsrcList() (list []uint32) {
	list = make([]uint32, rp.CsrcCount())
	for i := 0; i < len(list); i++ {
		list[i] = binary.BigEndian.Uint32(rp.buf[rtpHeaderLength+i*4:])
	}
	return
}

// Version returns the version as uint32 in host order.
func (pkt *RtpPacket) Payload() []byte {
	payOffset := int(pkt.CsrcCount()*4+rtpHeaderLength) + pkt.ExtensionLength()
	pad := 0
	if pkt.Padding() {
		pad = int(pkt.buf[pkt.inUse-1])
	}
	return pkt.buf[payOffset : pkt.inUse-pad]
}

const (
	// This field identifies the version of RTP. The version defined by this specification is two (2).
	versionOffset = 0
	// If the padding bit is set,
	// the packet contains one or more additional padding octets at the end which are not part of the payload.
	// P (Padding): (1 bit) Used to indicate if there are extra padding bytes at the end of the RTP packet.
	// A padding might be used to fill up a block of certain size, for example as required by an encryption algorithm.
	// The last byte of the padding contains the number of padding bytes that were added (including itself).[1
	paddingBitOffset = 0
	// If the extension bit is set,
	// the fixed header is followed by exactly one header extension,
	// with a format defined in Section 5.2.1.
	extensionBitOffset = 0
	// The CSRC count contains the number of CSRC identifiers that follow the fixed header.
	csrcCountOffset = 0
	// The interpretation of the marker is defined by a profile.
	// It is intended to allow significant events such as frame boundaries to be marked in the packet stream.
	markerOffset = 1
	// This field identifies the format of the RTP payload and determines its interpretation by the application.
	payloadTypeOffset = 1
	// The sequence number increments by one for each RTP data packet sent,
	// and may be used by the receiver to detect packet loss and to restore packet sequence.
	sequenceOffset = 2
	// 16 bits
	// The timestamp reflects the sampling instant of the first octet in the RTP data packet.
	// The sampling instant must be derived from a clock that increments monotonically and linearly in time to allow synchronization and jitter calculations
	timestampOffset = sequenceOffset + 2
	// 32 bits
	// The SSRC field identifies the synchronization source.
	ssrcOffsetRtp = timestampOffset + 4
	//  0 to 15 items, 32 bits each
	// The CSRC list identifies the contributing sources for the payload contained in this packet.
	// The number of identifiers is given by the CC field.
	// If there are more than 15 contributing sources, only 15 may be identified.
	// CSRC identifiers are inserted by mixers, using the SSRC identifiers of contributing sources.
	csrcOffsetRtp = ssrcOffsetRtp + 4
)

const (
	version2Bit  = 0x80
	extensionBit = 0x10
	paddingBit   = 0x20
	markerBit    = 0x80
	ccMask       = 0x0f
	ptMask       = 0x7f
	countMask    = 0x1f
)

/*
	该文件作为rtp解析、rtp构建包
	usage:
	n, remoteAddr, err := conn.ReadFromUDP(data)
	pkt := NewRtpPacket()
	pkt.Parse(data, n, remoteAddr)
	pkt.Print()
	// your process...
	go Process(pkt)
	// end of Process You MUST free pkt
	pkt.FreePacket()
*/

// Remote stores a remote addess in a transport independent way.
//
// The transport implementations construct UDP or TCP addresses and use them to send the data.
type Address struct {
	IPAddr   net.IP
	DataPort int
	CtrlPort int
	Zone     string
}

type RawPacket struct {
	inUse    int
	padTo    int
	conn     *net.UDPAddr
	fromAddr Address
	buf      []byte
}

type RtpPacket struct {
	RawPacket
	payloadLength int16
}

var freeRtpPacketList = make(chan *RtpPacket, freeRtpPacketListNum)

func NewRtpPacket() (pkt *RtpPacket) {
	select {
	case pkt = <-freeRtpPacketList: // Got one; nothing more to do.
	default:
		pkt = new(RtpPacket) // None free, so allocate a new one.
		pkt.buf = make([]byte, defaultMTULength)
	}
	// pkt.buffer[0] = version2Bit // RTP: V = 2, P, X, CC = 0
	pkt.inUse = rtpHeaderLength
	return pkt
}

func (pkt *RtpPacket) FreePacket() {

	pkt.buf[0] = 0 // invalidate RTP packet
	pkt.fromAddr.DataPort = 0
	pkt.fromAddr.IPAddr = nil
	pkt.inUse = 0

	select {
	case freeRtpPacketList <- pkt: // Packet on free list; nothing more to do.
	default: // Free list full, just carry on.
	}
}

// Print outputs a formatted dump of the RTP packet.
func (rp *RtpPacket) Print(label string) {
	fmt.Printf("RTP Packet at: %s\n", label)
	fmt.Printf("  fixed header dump:   %s\n", hex.EncodeToString(rp.buf[0:rtpHeaderLength]))
	fmt.Printf("    Version:           %d\n", (rp.buf[0]&0xc0)>>6)
	fmt.Printf("    Padding:           %t\n", rp.Padding())
	fmt.Printf("    Extension:         %t\n", rp.ExtensionBit())
	fmt.Printf("    Contributing SRCs: %d\n", rp.CsrcCount())
	fmt.Printf("    Marker:            %t\n", rp.Marker())
	fmt.Printf("    Payload type:      %d (0x%x)\n", rp.PayloadType(), rp.PayloadType())
	fmt.Printf("    Sequence number:   %d (0x%x)\n", rp.Sequence(), rp.Sequence())
	fmt.Printf("    Timestamp:         %d (0x%x)\n", rp.Timestamp(), rp.Timestamp())
	fmt.Printf("    SSRC:              %d (0x%x)\n", rp.Ssrc(), rp.Ssrc())

	if rp.CsrcCount() > 0 {
		cscr := rp.CsrcList()
		fmt.Printf("  CSRC list:\n")
		for i, v := range cscr {
			fmt.Printf("      %d: %d (0x%x)\n", i, v, v)
		}
	}
	if rp.ExtensionBit() {
		extLen := rp.ExtensionLength()
		fmt.Printf("  Extentsion length: %d\n", extLen)
		offsetExt := rtpHeaderLength + int(rp.CsrcCount()*4)
		fmt.Printf("    extension: %s\n", hex.EncodeToString(rp.buf[offsetExt:offsetExt+extLen]))
	}
	payOffset := rtpHeaderLength + int(rp.CsrcCount()*4) + rp.ExtensionLength()
	fmt.Printf("  payload: %s\n", hex.EncodeToString(rp.buf[payOffset:rp.inUse]))
}
