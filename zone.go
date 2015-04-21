package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/kteza1/ninjasphere-limitlessled/core"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/channels"
	"github.com/ninjasphere/go-ninja/model"
)

/*LimitlessLedBridge Struct for info about bridge.*/
type LimitlessLedBridge struct {
	*net.UDPConn
}

type LimitlessLedZone struct {
	driver            ninja.Driver
	info              *model.Device
	sendEvent         func(event string, payload interface{}) error
	onOffChannel      *channels.OnOffChannel
	brightnessChannel *channels.BrightnessChannel
	colorChannel      *channels.ColorChannel
	zoneNumber        int
	bridge            *LimitlessLedBridge
}

func newLimitlessLedZone(driver ninja.Driver, id int) *LimitlessLedZone {
	name := fmt.Sprintf("MiLight-Z%d", id)
	zone := &LimitlessLedZone{
		driver: driver,
		info: &model.Device{
			NaturalID:     fmt.Sprintf("light%d", id),
			NaturalIDType: "limitlessled zone",
			Name:          &name,
			Signatures: &map[string]string{
				"ninja:manufacturer": "LimitlessLed",
				"ninja:productName":  "LimitlessLed",
				"ninja:productType":  "Light",
				"ninja:thingType":    "light",
			},
		},
	}
	zone.onOffChannel = channels.NewOnOffChannel(zone)
	zone.brightnessChannel = channels.NewBrightnessChannel(zone)
	zone.colorChannel = channels.NewColorChannel(zone)
	return zone
}

/*SetOnOff -->*/
func (l *LimitlessLedZone) SetOnOff(state bool) error {
	zone := l.zoneNumber
	switch zone {
	case 1:
		if state == true {
			fmt.Println("Switch on zone1")
			l.bridge.SendCommand(core.ZONE1_ON)
		} else {
			fmt.Println("Switching off zone1")
			l.bridge.SendCommand(core.ZONE1_OFF)
		}
	case 2:
		if state == true {
			fmt.Println("Switch on zone2")
			l.bridge.SendCommand(core.ZONE2_ON)
		} else {
			fmt.Println("Switching off zone2")
			l.bridge.SendCommand(core.ZONE2_OFF)
		}
	case 3:
		if state == true {
			fmt.Println("Switch on zone3")
			l.bridge.SendCommand(core.ZONE3_ON)
		} else {
			fmt.Println("Switching off zone3")
			l.bridge.SendCommand(core.ZONE3_OFF)
		}
	case 4:
		if state == true {
			fmt.Println("Switch on zone4")
			l.bridge.SendCommand(core.ZONE4_ON)
		} else {
			fmt.Println("Switching off zone4")
			l.bridge.SendCommand(core.ZONE4_OFF)
		}
	}
	return nil
}

//ToggleOnOff -->
func (l *LimitlessLedZone) ToggleOnOff() error {
	log.Println("Toggling")
	return nil
}

func (l *LimitlessLedZone) SetColor(state *channels.ColorState) error {
	fmt.Println("setting color state to %v", state)
	return nil
}

func (l *LimitlessLedZone) SetBrightness(state float64) error {
	var brightness uint32 = (uint32)(state * 100)
	/* math to scale the range to [0, 27] */
	/* ref --> http://stackoverflow.com/questions/5294955/how-to-scale-down-a-range-of-numbers-with-a-known-min-and-max-value */
	brightness_byte := (byte)((27-0)*(brightness-0)/(100-0) + 0)
	c1 := byte(0x4E)
	c3 := byte(0x55)

	command := []byte{c1, brightness_byte, c3}

	zone := l.zoneNumber
	switch zone {
	case 1:
		l.bridge.SendCommand(core.ZONE1_ON)
		l.bridge.SendCommand(command)
	case 2:
		l.bridge.SendCommand(core.ZONE2_ON)
		l.bridge.SendCommand(command)
	case 3:
		l.bridge.SendCommand(core.ZONE3_ON)
		l.bridge.SendCommand(command)
	case 4:
		l.bridge.SendCommand(core.ZONE4_ON)
		l.bridge.SendCommand(command)
	}

	fmt.Printf("setting zone %d brightness to %d\n", zone, brightness)
	return nil
}

func (l *LimitlessLedZone) GetDeviceInfo() *model.Device {
	return l.info
}

func (l *LimitlessLedZone) GetDriver() ninja.Driver {
	return l.driver
}

//SetEventHandler -->
func (l *LimitlessLedZone) SetEventHandler(sendEvent func(event string, payload interface{}) error) {
	l.sendEvent = sendEvent
}

var reg, _ = regexp.Compile("[^a-z0-9]")

//SetName --> Exported by service/device schema
func (l *LimitlessLedZone) SetName(name *string) (*string, error) {
	log.Printf("Setting device name to %s", *name)
	safe := reg.ReplaceAllString(strings.ToLower(*name), "")
	if len(safe) > 16 {
		safe = safe[0:16]
	}
	log.Printf("Pretending we can only set 5 lowercase alphanum. Name now: %s", safe)
	l.sendEvent("renamed", safe)
	return &safe, nil
}

func Dial(host string) (*LimitlessLedBridge, error) {
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
	return &LimitlessLedBridge{s}, err
}

func (bridge *LimitlessLedBridge) SendCommand(command []byte) {
	fmt.Println("Sending command")
	/* Sending each command twice since bridge can get slow sometimes */
	for i := 0; i < 2; i++ {
		_, err := bridge.Write(command)
		if err != nil {
			fmt.Println("Error writing")
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
	time.Sleep(time.Millisecond * 150)
}
