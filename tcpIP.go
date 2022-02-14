package main

import (
	"fmt"
	"net"
	"time"
)

const (
	serverPort = ":34933"
	serverIP   = "10.100.23.240"
	localIP    = "10.100.23.156"
)

func main() {

	con, err := net.Dial("tcp", serverIP+serverPort)
	errorHandler(err)

	go receiveMessages(con)
	con.Write([]byte("Connect to 10.100.23.156:34933\x00"))
	//con.Write([]byte("Hello worldlies\x00"))

	defer con.Close()

	time.Sleep(time.Second * 20)

}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func receiveMessages(con net.Conn) {
	buf := make([]byte, 1024)
	for {
		length, err := con.Read(buf)
		errorHandler(err)
		fmt.Printf("Received: %s\r\n", buf[:length])
	}
}
