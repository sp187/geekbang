package decoder

import (
	"encoding/binary"
	"net"
)

type GoimMessage struct {
	PackageLen int
	HeaderLen  int
	Version    int
	Operation  int
	Seq        uint
	Header     []byte
	Body       []byte
}

func GoimDecoder(conn net.Conn) GoimMessage {
	msg := GoimMessage{}
	pkgLenBytes := make([]byte, 4)
	conn.Read(pkgLenBytes)
	// 读取包长度
	pkgLen := binary.BigEndian.Uint32(pkgLenBytes)
	msg.PackageLen = int(pkgLen)
	// 读取包内容
	msgBytes := make([]byte, pkgLen-4)
	conn.Read(msgBytes)

	headerLen := binary.BigEndian.Uint16(msgBytes[:2])
	msg.HeaderLen = int(headerLen)

	version := binary.BigEndian.Uint16(msgBytes[2:4])
	msg.Version = int(version)

	operation := binary.BigEndian.Uint32(msgBytes[4:8])
	msg.Operation = int(operation)

	seqId := binary.BigEndian.Uint32(msgBytes[8:12])
	msg.Seq = uint(seqId)

	msg.Header = msgBytes[12 : 12+headerLen]
	msg.Body = msgBytes[12+headerLen:]
	return msg
}
