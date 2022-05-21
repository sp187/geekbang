package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	conn, err := net.Dial("tcp", "localhost:10001")
	if err != nil {
		panic("connect to server error")
	}
	defer conn.Close()
	tic := time.NewTicker(3 * time.Second)
	defer tic.Stop()
	for {
		select {
		case <-tic.C:
			msg := makeMessage()
			fmt.Printf("msg len :%d\n", len(msg))
			_, err := conn.Write(msg)
			if err != nil {
				return
			}
		}
	}
}

func makeMessage() []byte {
	bodyLen := rand.Int() % 20
	if bodyLen == 0 {
		bodyLen++
	}
	bodyStr := ""
	for i := 0; i < bodyLen; i++ {
		bodyStr += "a"
	}
	msgLen := bodyLen + 2
	msg := make([]byte, msgLen)
	binary.BigEndian.PutUint16(msg[:2], uint16(bodyLen))
	copy(msg[2:], bodyStr)
	return msg
}
