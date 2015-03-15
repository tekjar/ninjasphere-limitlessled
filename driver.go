package main

import (
	"fmt"
	"log"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/events"
	"github.com/ninjasphere/go-ninja/support"
)

var info = ninja.LoadModuleInfo("./package.json")

/*LimitlessLedDriver --> Struct for LimitlessLed driver.*/
type LimitlessLedDriver struct {
	support.DriverSupport
	config *LimitlessLedDriverConfig
}

/*LimitlessLedDriverConfig --> Struct for LimitlessLed driver configuration.*/
type LimitlessLedDriverConfig struct {
	Initialised     bool
	NumberOfBridges int
}

func defaultConfig() *LimitlessLedDriverConfig {
	return &LimitlessLedDriverConfig{
		Initialised:     false,
		NumberOfBridges: 1,
	}
}

/*NewLimitlessLedDriver --> initializes a new LimitlessLed Driver.*/
func NewLimitlessLedDriver() (*LimitlessLedDriver, error) {
	driver := &LimitlessLedDriver{}
	err := driver.Init(info)
	if err != nil {
		log.Fatalf("Failed to initialize Orvibo driver: %s", err)
	}
	//exposes the driver to ninja sphere framework
	err = driver.Export(driver)
	if err != nil {
		log.Fatalf("Failed to export LimitlessLed driver: %s", err)
	}
	return driver, nil
}

/*OnPairingRequest --> */
func (d *LimitlessLedDriver) OnPairingRequest(pairingRequest *events.PairingRequest, values map[string]string) bool {
	log.Printf("RTR. Pairing request received from %s for %d seconds", values["deviceId"], pairingRequest.Duration)
	return true
}

/*Start -->  */
func (d *LimitlessLedDriver) Start(config *LimitlessLedDriverConfig) error {
	log.Printf("Driver Starting with config %v", config)
	bridgeIps := [4]string{"192.168.0.100:8899", "192.168.0.101:8899", "192.168.0.102:8899", "192.168.0.103:8899"}
	d.config = config
	if !d.config.Initialised {
		d.config = defaultConfig()
	}
	/* Don't let it cross more than 4 for now */
	for i := 0; i < d.config.NumberOfBridges; i++ {
		go func() {
			device := NewLimitlessLedBridge(d, i, bridgeIps[i])
			_, err := device.Dial(bridgeIps[i])
			if err != nil {
				fmt.Println("Something wrong")
				return
			}
			/* If Dail is successful, expose device and channels */
			err = d.Conn.ExportDevice(device)
			if err != nil {
				log.Fatalf("Failed to export the bridge %d: %s", i, err)
			}
			/* Each bridge is configured with 4 zones. zone = on-off channel coz control is through
			zones. But a zone can have multiple lights configured */
			err = d.Conn.ExportChannel(device, device.onOffChannel1, "on-off")
			if err != nil {
				log.Fatalf("Failed to export bridge's zone1 on off channel %d: %s", i, err)
			}
			err = d.Conn.ExportChannel(device, device.onOffChannel2, "on-off")
			if err != nil {
				log.Fatalf("Failed to export bridge's zone2 on off channel %d: %s", i, err)
			}
			err = d.Conn.ExportChannel(device, device.onOffChannel3, "on-off")
			if err != nil {
				log.Fatalf("Failed to export bridge's zone3 on off channel %d: %s", i, err)
			}
			err = d.Conn.ExportChannel(device, device.onOffChannel4, "on-off")
			if err != nil {
				log.Fatalf("Failed to export bridge's zone4 on off channel %d: %s", i, err)
			}
		}()
	}

	return d.SendEvent("config", config)
}

//In -->
type In struct {
	Name string
}

//Out -->
type Out struct {
	Age  int
	Name string
}

//Blarg -->
func (d *LimitlessLedDriver) Blarg(in *In) (*Out, error) {
	log.Printf("GOT INCOMING! %s", in.Name)
	return &Out{
		Name: in.Name,
		Age:  30,
	}, nil
}
