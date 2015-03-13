package main

import (
	"fmt"
	"github.com/kteza1/ninjasphere-limitlessled/core"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/channels"
	"github.com/ninjasphere/go-ninja/model"
	"log"
	"regexp"
	"strings"
)

// Struct for info about bridge
type LimitlessLedBridge struct {
	driver        ninja.Driver
	info          *model.Device
	sendEvent     func(event string, payload interface{}) error
	onOffChannel1 *channels.OnOffChannel
	onOffChannel2 *channels.OnOffChannel
	onOffChannel3 *channels.OnOffChannel
	onOffChannel4 *channels.OnOffChannel
	bridge        core.Bridge
}

func NewLimitlessLedBridge(driver ninja.Driver, id core.Bridge) *LimitlessLedBridge {
	name := id.Name
	socket := &LimitlessLedBridge{
		driver: driver,
		bridge: id,
		info: &model.Device{
			NaturalID:     fmt.Sprintf("socket%d", 1122334455),
			NaturalIDType: "socket",
			Name:          &name,
			Signatures: &map[string]string{
				"ninja:manufacturer": "LimitlessLed",
				"ninja:productName":  "LimitlessLed",
				"ninja:productType":  "Socket",
				"ninja:thingType":    "socket",
			},
		},
	}
	socket.onOffChannel = channels.NewOnOffChannel(socket)
	return socket
}

func (l *LimitlessLedBridge) GetDeviceInfo() *model.Device {
	return l.info
}
func (l *LimitlessLedBridge) GetDriver() ninja.Driver {
	return l.driver
}
func (l *LimitlessLedBridge) SetOnOff(state bool) error {
	allone.SetState(state, l.Socket.MACAddress)
	return nil
}

func (l *LimitlessLedBridge) SetEventHandler(sendEvent func(event string, payload interface{}) error) {
	l.sendEvent = sendEvent
}

var reg, _ = regexp.Compile("[^a-z0-9]")

// Exported by service/device schema
func (l *LimitlessLedBridge) SetName(name *string) (*string, error) {
	log.Printf("Setting device name to %s", *name)
	safe := reg.ReplaceAllString(strings.ToLower(*name), "")
	if len(safe) > 16 {
		safe = safe[0:16]
	}
	log.Printf("Pretending we can only set 5 lowercase alphanum. Name now: %s", safe)
	l.Socket.Name = safe
	l.sendEvent("renamed", safe)
	return &safe, nil
}
