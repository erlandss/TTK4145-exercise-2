package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":20017")
	fmt.Println(err)
	pc, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("listenpacket<s error")
		panic(err)
	}
	defer pc.Close()

	go readServer(*pc)
	go writeToServer()
	time.Sleep(20 * time.Second)
}

func readServer(pc net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		n, addr, err := pc.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("readfrom error")
			panic(err)
		}

		fmt.Printf("%s sent this: %s\n", addr, buf[:n])
	}
}

func writeToServer() {
	con, err := net.Dial("udp", "10.100.23.240:20017")
	errorHandler(err)

	defer con.Close()

	for {
		_, err := con.Write([]byte("hei"))
		errorHandler(err)
		time.Sleep(time.Second)
	}

}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
