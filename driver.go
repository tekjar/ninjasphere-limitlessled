package main

import (
	"fmt"
	"log"

	"github.com/kteza1/ninjasphere-limitlessled/core"
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
		NumberOfBridges: 1, //No. of lights ?? Does it matter ??
	}
}

/*NewLimitlessLedDriver --> initializes a new LimitlessLed Driver.*/
func NewLimitlessLedDriver() (*LimitlessLedDriver, error) {
	driver := &LimitlessLedDriver{}
	err := driver.Init(info)
	if err != nil {
		log.Fatalf("Failed to initialize LimitlessLed driver: %s", err)
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
	fmt.Println("RTR. Pairing request received from %s for %d seconds", values["deviceId"], pairingRequest.Duration)
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
	var id int = 0
	/* Don't let it cross more than 4 for now */
	for i := 0; i < d.config.NumberOfBridges; i++ {
		fmt.Println("Creating connection to", bridgeIps[i])
		bridge, err := Dial(bridgeIps[i])
		if err != nil {
			fmt.Println("Something wrong while trying to connect to bridge")
			return err
		}
		/* Switch all off just for the confirmation that the bridge is connected */
		bridge.SendCommand(core.ALL_OFF)

		/* zone 1*/
		id++
		device1 := newLimitlessLedZone(d, id) /* A new id for each zone */
		device1.bridge = bridge
		device1.zoneNumber = 1
		/* If Dail is successful, expose zones*/
		err = d.Conn.ExportDevice(device1)
		if err != nil {
			log.Printf("Failed to export zone 1. id =  %d, err = %s", id, err)
		}
		err = d.Conn.ExportChannel(device1, device1.onOffChannel, "on-off")
		if err != nil {
			log.Printf("Failed to export zone 1 on off channel. err = %s", err)
		}

		/* zone2 */
		id++
		device2 := newLimitlessLedZone(d, id)
		device2.bridge = bridge
		device2.zoneNumber = 2
		/* If Dail is successful, expose zones*/
		err = d.Conn.ExportDevice(device2)
		if err != nil {
			log.Printf("Failed to export zone 2. id =  %d, err = %s", id, err)
		}
		err = d.Conn.ExportChannel(device2, device2.onOffChannel, "on-off")
		if err != nil {
			log.Printf("Failed to export zone 2 on off channel. err = %s", err)
		}

		/* zone3 */
		id++
		device3 := newLimitlessLedZone(d, id)
		device3.bridge = bridge
		device3.zoneNumber = 3
		/* If Dail is successful, expose zones*/
		err = d.Conn.ExportDevice(device3)
		if err != nil {
			log.Printf("Failed to export zone 3. id =  %d, err = %s", id, err)
		}
		err = d.Conn.ExportChannel(device3, device3.onOffChannel, "on-off")
		if err != nil {
			log.Printf("Failed to export zone 3 on off channel. err = %s", err)
		}

		/* zone4 */
		id++
		device4 := newLimitlessLedZone(d, id)
		device4.bridge = bridge
		device4.zoneNumber = 4
		/* If Dail is successful, expose zones*/
		err = d.Conn.ExportDevice(device4)
		if err != nil {
			log.Printf("Failed to export zone 4. id =  %d, err = %s", id, err)
		}
		err = d.Conn.ExportChannel(device4, device4.onOffChannel, "on-off")
		if err != nil {
			log.Printf("Failed to export zone 4 on off channel. err = %s", err)
		}

	}

	return d.SendEvent("config", config)
}
