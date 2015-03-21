package main

import (
	"fmt"
	"net"
	"time"
)

var (
	ALL_ON          = []byte{0x42, 0x00, 0x55}
	ALL_WHITE       = []byte{0xC2, 0x00, 0x55}
	ALL_OFF         = []byte{0x41, 0x00, 0x55}
	ALL_DISCO       = []byte{0x4D, 0x00, 0x55}
	BRIGHTNESS_UP   = []byte{0x3C, 0x00, 0x55}
	BRIGHTNESS_DOWN = []byte{0x34, 0x00, 0x55}
	WARMER          = []byte{0x3E, 0x00, 0x55}
	COOLER          = []byte{0x3F, 0x00, 0x55}
	ALL_ON_FULL     = []byte{0xB5, 0x00, 0x55}
	ALL_NIGHTLIGHT  = []byte{0xB9, 0x00, 0x55}
	ZONE1_ON        = []byte{0x45, 0x00, 0x55}
	ZONE1_OFF       = []byte{0x46, 0x00, 0x55}
	ZONE2_ON        = []byte{0x47, 0x00, 0x55}
	ZONE2_OFF       = []byte{0x48, 0x00, 0x55}
	ZONE3_ON        = []byte{0x49, 0x00, 0x55}
	ZONE3_OFF       = []byte{0x4A, 0x00, 0x55}
	ZONE4_ON        = []byte{0x4B, 0x00, 0x55}
	ZONE4_OFF       = []byte{0x4C, 0x00, 0x55}
)

type Bridge struct {
	*net.UDPConn
}

func Dial(host string) (*Bridge, error) {
	addr, err := net.ResolveUDPAddr("udp4", host)

	if err != nil {
		fmt.Println("Error Resolving")
		return nil, err
	}

	s, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		fmt.Println("Error dialing")
		return nil, err
	}
	return &Bridge{s}, err
}

func (bridge *Bridge) SendCommand(command []byte) {
	fmt.Println("Sending command")
	_, err := bridge.Write(command)
	if err != nil {
		fmt.Println("Error writing")
		return
	}
	time.Sleep(time.Millisecond * 100)
}
