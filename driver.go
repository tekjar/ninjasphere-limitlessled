package main

import (
	"fmt"
	"github.com/kteza1/ninjasphere-limitlessled/core"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/support"
	"log"
	"time"
)

var info = ninja.LoadModuleInfo("./package.json")
var ready = false

type LimitlessLedDriver struct {
	support.DriverSupport
	config *LimitlessLedDriverConfig
}
type LimitlessLedDriverConfig struct {
	Initialised   bool
	NumberOfZones int
}

func defaultConfig() *OrviboDriverConfig {
	return &OrviboDriverConfig{
		Initialised:   false,
		NumberOfZones: 4,
	}
}
func NewDriver() (*LimitlessLedDriver, error) {
	driver := &LimitlessLedDriver{}
	err := driver.Init(info)
	if err != nil {
		log.Fatalf("Failed to initialize Orvibo driver: %s", err)
	}
	err = driver.Export(driver)
	if err != nil {
		log.Fatalf("Failed to export Orvibo driver: %s", err)
		allone.Close()
	}
	return driver, nil
}

/* Driver will see core's Events channel to perform actions */
func (d *LimitlessLedDriver) Start(config *LimitlessLedDriverConfig) error {
	log.Printf("Driver Starting with config %v", config)
	d.config = config
	if !d.config.Initialised {
		d.config = defaultConfig()
	}
	var firstDiscover, firstSubscribe, firstQuery, autoDiscover, resubscribe chan bool
	var device *OrviboSocket
	d.SendEvent("config", config)
	go func() {
		device = NewOrviboSocket(d, msg.SocketInfo)
		device.Socket.Name = msg.Name
		err := d.Conn.ExportDevice(device)
		err = d.Conn.ExportChannel(device, device.onOffChannel, "on-off")
		if err != nil {
			log.Fatalf("Failed to export Orvibo socket on off channel %s: %s", msg.SocketInfo.MACAddress, err)
			allone.Close()
		}

		allone.CheckForMessages()
		for {
			allone.CheckForMessages()
			select {
			case msg := <-allone.Events:
				fmt.Println("!!!T Type:", msg.Name)
				switch msg.Name {
				case "ready":
				case "subscribed":
				case "queried":
					fmt.Println("We've queried. Name is:", msg.SocketInfo.Name)
					device = NewOrviboSocket(d, msg.SocketInfo)
					device.Socket.Name = msg.Name
					err := d.Conn.ExportDevice(device)
					err = d.Conn.ExportChannel(device, device.onOffChannel, "on-off")
					if err != nil {
						log.Fatalf("Failed to export Orvibo socket on off channel %s: %s", msg.SocketInfo.MACAddress, err)
						allone.Close()
					}
					allone.Discover()
					allone.Subscribe()
					allone.CheckForMessages()
				case "statechanged":
					fmt.Println("State changed to:", msg.SocketInfo.State)
					allone.CheckForMessages()
				}
			}
			allone.CheckForMessages()
		}
	}()
	return d.SendEvent("config", config)
}
func (d *OrviboSocket) Stop() error {
	allone.Close()
	return fmt.Errorf("This driver does not support being stopped. YOU HAVE NO POWER HERE.")
}

type In struct {
	Name string
}
type Out struct {
	Age  int
	Name string
}

func (d *OrviboDriver) Blarg(in *In) (*Out, error) {
	log.Printf("GOT INCOMING! %s", in.Name)
	return &Out{
		Name: in.Name,
		Age:  30,
	}, nil
}
func setInterval(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)
	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()
	return stop
}
