package core

import

// For outputting stuff
(
	"fmt"
	"net"
	"time"
) // For networking stuff
// For exiting
// For reversing strings

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
	ZONE2_SYNCPAIR  = []byte{0x47, 0x00, 0x55}
	ZONE2_OFF       = []byte{0x33, 0x00, 0x55}
)

//Bridge LimitlessLed bridge
type Bridge struct {
	*net.UDPConn        //extending bridge capabilities
	ip           string // The name of the bridge
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
	return &Bridge{s, host}, err
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
