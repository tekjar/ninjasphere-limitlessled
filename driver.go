package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

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

type Milight struct {
	XMLName xml.Name `xml:"Milight"`
	Bridge  struct {
		XMLName xml.Name `xml:"Bridge"`
		IP      string   `xml:"ip,attr"`
		Zone1   []string `xml:"zone1"`
		Zone2   []string `xml:"zone2"`
		Zone3   []string `xml:"zone3"`
		Zone4   []string `xml:"zone4"`
	}
}

type Zone struct {
	Value string `xml:",chardata"`
	Inst  int    `xml:"inst,attr"`
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
	miXML, err := os.Open("milights.xml")
	if err != nil {
		fmt.Println("Error opening the file")
		log.Fatalf("Unable to find the milight xml database")
	}

	defer miXML.Close()

	miData, _ := ioutil.ReadAll(miXML)
	var xMilight Milight
	xml.Unmarshal(miData, &xMilight)
	bridgeIP := xMilight.Bridge.IP
	d.config = config
	if !d.config.Initialised {
		d.config = defaultConfig()
	}
	var id int = 0
	/* This is 1 for now. Will be taken from config file later */
	for i := 0; i < d.config.NumberOfBridges; i++ {
		fmt.Println("Creating connection to", bridgeIP)
		bridge, err := Dial(bridgeIP)
		if err != nil {
			fmt.Println("Something wrong while trying to connect to bridge")
			return err
		}
		/* Blink all Lights. Confirmation that the driver is ready*/
		time.Sleep(time.Millisecond * 150)
		bridge.SendCommand(core.ALL_OFF)
		time.Sleep(time.Millisecond * 150)
		bridge.SendCommand(core.ALL_ON)
		time.Sleep(time.Millisecond * 150)
		bridge.SendCommand(core.ALL_OFF)

		/* zone 1*/
		id++
		device1 := newLimitlessLedZone(d, id) /* A new id for each zone */
		device1.bridge = bridge
		device1.zoneNumber = 1
		/* If Dail is successful, expose Device*/
		err = d.Conn.ExportDevice(device1)
		if err != nil {
			log.Printf("Failed to export zone 1. id =  %d, err = %s", id, err)
		}
		/* Start exporting channels for this device */
		err = d.Conn.ExportChannel(device1, device1.onOffChannel, "on-off")
		if err != nil {
			log.Printf("Failed to export zone 1 on off channel. err = %s", err)
		}

		err = d.Conn.ExportChannel(device1, device1.brightnessChannel, "brightness")
		if err != nil {
			log.Fatalf("Failed to export fake light brightness channel %d: %s", i, err)
		}
		err = d.Conn.ExportChannel(device1, device1.colorChannel, "color")
		if err != nil {
			log.Fatalf("Failed to export fake color channel %d: %s", i, err)
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
		/* Export channels */
		err = d.Conn.ExportChannel(device2, device2.onOffChannel, "on-off")
		if err != nil {
			log.Printf("Failed to export zone 2 on off channel. err = %s", err)
		}
		err = d.Conn.ExportChannel(device2, device2.brightnessChannel, "brightness")
		if err != nil {
			log.Fatalf("Failed to export fake light brightness channel %d: %s", i, err)
		}
		err = d.Conn.ExportChannel(device2, device2.colorChannel, "color")
		if err != nil {
			log.Fatalf("Failed to export fake color channel %d: %s", i, err)
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
		err = d.Conn.ExportChannel(device3, device3.brightnessChannel, "brightness")
		if err != nil {
			log.Fatalf("Failed to export fake light brightness channel %d: %s", i, err)
		}
		err = d.Conn.ExportChannel(device3, device3.colorChannel, "color")
		if err != nil {
			log.Fatalf("Failed to export fake color channel %d: %s", i, err)
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
		err = d.Conn.ExportChannel(device4, device4.brightnessChannel, "brightness")
		if err != nil {
			log.Fatalf("Failed to export fake light brightness channel %d: %s", i, err)
		}
		err = d.Conn.ExportChannel(device4, device4.colorChannel, "color")
		if err != nil {
			log.Fatalf("Failed to export fake color channel %d: %s", i, err)
		}

	}

	return d.SendEvent("config", config)
}
