package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

type TypeTaggedJSON struct {
	TypeId string
	JSON   []byte
}

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down MotorDirection = -1
	MD_Stop MotorDirection = 0
)

type ElevatorBehavior int

const (
	EB_Idle ElevatorBehavior = iota
	EB_DoorOpen
	EB_Moving
)

type ClearRequestVariant int

const (
	CV_All ClearRequestVariant = iota
	CV_InDirn
)

type Config struct {
	ClearRequestVariant ClearRequestVariant
	DoorOpenDuration_s  int64 //
}

type Elevator struct {
	Id        int
	Floor     int
	Dirn      MotorDirection
	Requests  [3][4]bool
	Behavior  ElevatorBehavior
	Config    Config
	Operative bool
}

type OrderState int

const (
	Order_Unknown OrderState = iota
	Order_None
	Order_Unconfirmed
	Order_Confirmed
)

type NetworkMessage struct {
	Id          int
	ElevStates  [3]Elevator
	OrderStates [3][4][3]OrderState
}

func main() {
	//addr, err := net.ResolveUDPAddr("udp", ":10008")

	go readServer()
	//go writeToServer()
	time.Sleep(20 * time.Second)
	for {

	}
}

func DialBroadcastUDP(port int) net.PacketConn {
	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		fmt.Println("Error: Socket:", err)
	}
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		fmt.Println("Error: SetSockOpt REUSEADDR:", err)
	}
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	if err != nil {
		fmt.Println("Error: SetSockOpt BROADCAST:", err)
	}
	syscall.Bind(s, &syscall.SockaddrInet4{Port: port})
	if err != nil {
		fmt.Println("Error: Bind:", err)
	}

	f := os.NewFile(uintptr(s), "")
	conn, err := net.FilePacketConn(f)
	if err != nil {
		fmt.Println("Error: FilePacketConn:", err)
	}
	f.Close()

	return conn
}

func readServer() {
	conn := DialBroadcastUDP(10013)
	buff := make([]byte, 2048)
	message := new(NetworkMessage)
	for {
		n, _, _ := conn.ReadFrom(buff)
		var ttj TypeTaggedJSON
		json.Unmarshal(buff[0:n], &ttj)
		message.Id = -1
		json.Unmarshal(ttj.JSON, message)

		fmt.Printf("%d sent this: %v\n", message.Id, message.OrderStates[0][3])
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
