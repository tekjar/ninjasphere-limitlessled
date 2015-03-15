package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/channels"
	"github.com/ninjasphere/go-ninja/model"
)

/*LimitlessLedBridge Struct for info about bridge.*/
type LimitlessLedBridge struct {
	driver        ninja.Driver
	info          *model.Device
	sendEvent     func(event string, payload interface{}) error
	onOffChannel1 *channels.OnOffChannel
	onOffChannel2 *channels.OnOffChannel
	onOffChannel3 *channels.OnOffChannel
	onOffChannel4 *channels.OnOffChannel
	ipPort        string
	conn          *net.UDPConn
	currentZone   int //for SetOnOff() method on the relevant channel based on current zone
}

/*NewLimitlessLedBridge function initializes a new bridge. */
func NewLimitlessLedBridge(driver ninja.Driver, id int, ipPort string) *LimitlessLedBridge {
	name := fmt.Sprintf("LimitlessLedBridge %d", id)
	bridge := &LimitlessLedBridge{
		driver: driver,
		info: &model.Device{
			NaturalID:     fmt.Sprintf("socket%d", 1122334455),
			NaturalIDType: "limitlessled bridge",
			Name:          &name,
			Signatures: &map[string]string{
				"ninja:manufacturer": "LimitlessLed",
				"ninja:productName":  "LimitlessLed",
				"ninja:productType":  "bridge",
				"ninja:thingType":    "light",
			},
		},
		ipPort: ipPort,
	}
	bridge.onOffChannel1 = channels.NewOnOffChannel(bridge)
	bridge.onOffChannel2 = channels.NewOnOffChannel(bridge)
	bridge.onOffChannel3 = channels.NewOnOffChannel(bridge)
	bridge.onOffChannel4 = channels.NewOnOffChannel(bridge)
	return bridge
}

/*GetDeviceInfo --> Function for getting bridge info */
func (l *LimitlessLedBridge) GetDeviceInfo() *model.Device {
	return l.info
}

/*GetDriver -->.*/
func (l *LimitlessLedBridge) GetDriver() ninja.Driver {
	return l.driver
}

/*SetOnOff -->*/
func (l *LimitlessLedBridge) SetOnOff(state bool) error {
	zone := l.currentZone
	switch zone {
	case 1:
		if state == true {
			fmt.Println("Switch on zone1")
		} else {
			fmt.Println("Switching off zone1")
		}
	case 2:
		if state == true {
			fmt.Println("Switch on zone2")
		} else {
			fmt.Println("Switching off zone2")
		}
	case 3:
		if state == true {
			fmt.Println("Switch on zone3")
		} else {
			fmt.Println("Switching off zone3")
		}
	case 4:
		if state == true {
			fmt.Println("Switch on zone4")
		} else {
			fmt.Println("Switching off zone4")
		}
	}
	return nil
}

//ToggleOnOff -->
func (l *LimitlessLedBridge) ToggleOnOff() error {
	log.Println("Toggling")
	return nil
}

//SetEventHandler -->
func (l *LimitlessLedBridge) SetEventHandler(sendEvent func(event string, payload interface{}) error) {
	l.sendEvent = sendEvent
}

var reg, _ = regexp.Compile("[^a-z0-9]")

//SetName --> Exported by service/device schema
func (l *LimitlessLedBridge) SetName(name *string) (*string, error) {
	log.Printf("Setting device name to %s", *name)
	safe := reg.ReplaceAllString(strings.ToLower(*name), "")
	if len(safe) > 16 {
		safe = safe[0:16]
	}
	log.Printf("Pretending we can only set 5 lowercase alphanum. Name now: %s", safe)
	l.sendEvent("renamed", safe)
	return &safe, nil
}

func (l *LimitlessLedBridge) Dial(host string) (*LimitlessLedBridge, error) {
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
	l.conn = s
	return l, err
}

func (bridge *LimitlessLedBridge) SendCommand(command []byte) {
	fmt.Println("Sending command")
	_, err := bridge.conn.Write(command)
	if err != nil {
		fmt.Println("Error writing")
		return
	}
	time.Sleep(time.Millisecond * 100)
}
