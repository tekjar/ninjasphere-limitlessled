package core

import (
	"encoding/hex"
	"fmt"     // For outputting stuff
	"net"     // For networking stuff
	"os"      // For exiting
	"strings" // For reversing strings
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
	ZONE2_SYNCPAIR  = []byte{0x47, 0x00, 0x55}
	ZONE2_OFF       = []byte{0x33, 0x00, 0x55}
)

// LimitlessLed bridge
type Bridge struct {
	*net.UDPConn        //extending bridge capabilities
	Name         string // The name of the socket (e.g. "Christmas Lights")
	ip_port      string // 192.168.0.100:8899
}

// EventStruct is what we pass back to our calling code via channels
type EventStruct struct {
	Name   string // The name of the event (e.g. ready, socketfound)
	bridge Bridge // And our Socket struct so we can look at IP address, MAC etc.
}

var conn *net.UDPConn // Our UDP connection, read and write
var msg []byte

var Events = make(chan EventStruct)

// SphereListner gets our UDP socket ready for reading and writing.
func SphereListner() {
	if getLocalIP() == "" {
		fmt.Println("Error: Can't determine local IP address. Exiting!")
		os.Exit(1)
	} else {
		fmt.Println("Local IP is:", getLocalIP())
	}
	udpAddr, err := net.ResolveUDPAddr("udp4", ":10000") // Get our address ready for listening
	if err != nil {
		fmt.Println("Resolve:", err)
		os.Exit(1)
	}
	conn, err = net.ListenUDP("udp", udpAddr) // Now we listen on the address we just resolved
	if err != nil {
		fmt.Println("Listen:", err)
		os.Exit(1)
	}
	go func() { Events <- EventStruct{"ready", Bridge{}} }()
}

/* Sphere's check for messages */
func CheckForMessages() { // Now we're checking for messages
	var buf [1024]byte // We want to get 1024 bytes of messages (is this enough? Need to check!)
	go func() {        // Rading from UDP blocks
		n, addr, _ := conn.ReadFromUDP(buf[0:])
		msg = buf[0:n]                                 // Set this property so other functions can use it. n is how many bytes we grabbed from UDP
		if n > 0 && addr.IP.String() != getLocalIP() { // If we've got more than 0 bytes
			fmt.Println("Yo, Message was found:", n)
			go handleMessage(hex.EncodeToString(msg), addr) // Hand it off to our handleMessage func. We pass on the message and the address (for replying to messages)
		}
		msg = nil // Clear out our msg property so we don't run handleMessage on old data
	}() // Read from UDP connection. [0:] is slice stuff that says "shove everything in the first section of the byte and go until we've extracted all data"
}

// SetState sets the state of a socket, given its MAC address
func SetState(state bool, macAdd string) {
	var statebit string
	if state == true {
		statebit = "01"
	} else {
		statebit = "00"
	}
	sendMessage("686400176463"+macAdd+twenties+"00000000"+statebit, sockets[macAdd].IP)
	go func() { Events <- EventStruct{"stateset", *sockets[macAdd]} }()
}

//send message to LimitlessLed Bridge
func sendMessage(msg string, ip *net.UDPAddr) {
	fmt.Println("Sending message:", msg, "to", ip.String())
	// Turn this hex string into bytes for sending
	buf, _ := hex.DecodeString(msg)
	// Resolve our address, ready for sending data
	udpAddr, _ := net.ResolveUDPAddr("udp4", ip.String())
	// Actually write the data and send it off
	// _ lets us ignore "declared but not used" errors. If we replace _ with n,
	// We'd have to use n somewhere (e.g. fmt.Println(n)), but _ lets us ignore that
	_, _ = conn.WriteToUDP(buf, udpAddr)
	// If we've got an error
	return
}

func handleMessage(message string, addr *net.UDPAddr) {
	if len(message) == 0 {
		return
	}
	commandID := message[8:12] // What command we've received back
	macStart := strings.Index(message, "accf")
	macAdd := message[macStart:(macStart + 12)] // The MAC address of the socket responding
	fmt.Println("Message:", message, "IP:", addr.IP.String(), "MAC:", macAdd, "CID:", commandID, "Time:", time.Now())
	switch commandID {
	case "all_on":
		go func() { Events <- EventStruct{"socketfound", *sockets[macAdd]} }()
	case "all_off":
		go func() { Events <- EventStruct{"subscribed", *sockets[macAdd]} }()
	case "all_disco": // We've queried our socket, this is the data back
		Events <- EventStruct{"queried", *sockets[macAdd]}
	case "all_white":
		Events <- EventStruct{"statechanged", *sockets[macAdd]}
	default:
		fmt.Println("Wrong option. Provide some thing appropriate")
	}

}

func Close() bool {
	err := conn.Close()
	if err != nil {
		fmt.Println("Error closing socket:", err)
		return false
	}
	return true
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Oops: " + err.Error() + "\n")
		return ""
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPAddr:
				return v.IP.String()
			}
		}
	}
	return ""
}
