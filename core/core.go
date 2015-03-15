package core

import

// For outputting stuff

"net"

// For networking stuff
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
