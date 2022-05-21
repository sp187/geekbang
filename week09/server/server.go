package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:10001")
	if err != nil {
		panic("listen address error")
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept error %v\n", err)
			continue
		}
		//go fixLenHandler(conn)
		//go delimiterBasedHandler(conn)
		go protocolHandler(conn)
	}
}

// 固定长度的解包方式，需要客户端每次发送固定长度的信息，可用于简单的消息传递，如心跳信息。
func fixLenHandler(conn net.Conn) {
	defer conn.Close()
	Length := 10

	for {
		b := make([]byte, Length)
		n, err := conn.Read(b)
		if err != nil {
			fmt.Printf("read error %v\n", err)
			return
		}
		if n > 0 {
			fmt.Println(string(b))
		}
	}

}

// 按约定的分隔符进行解包。
func delimiterBasedHandler(conn net.Conn) {
	defer conn.Close()
	Delimiter := byte('&')
	rd := bufio.NewReader(conn)
	for {
		line, err := rd.ReadSlice(Delimiter)
		if err != nil {
			fmt.Printf("read error %v\n", err)
			return
		}
		fmt.Println(string(line))
	}
}

// 自定义协议进行解包，信息分为固定的协议头和不定长的信息体组成。
func protocolHandler(conn net.Conn) {
	defer conn.Close()
	headerLen := 2
	for {
		headerByte := make([]byte, headerLen)
		_, err := conn.Read(headerByte)
		if err != nil {
			fmt.Printf("read error %v\n", err)
			return
		}
		bodyLen := binary.BigEndian.Uint16(headerByte)
		fmt.Printf("body len: %d\n", bodyLen)
		body := make([]byte, bodyLen)
		n, err := conn.Read(body)
		if err != nil {
			fmt.Printf("read error %v\n", err)
			return
		}
		if n != int(bodyLen) {
			fmt.Printf("invalid body len %d\n", n)
			return
		}
		fmt.Println(string(body))
	}
}
