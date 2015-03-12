package main

import (
	"fmt"
	"os"
)

func main() {
	/* At initial bridge boot, router seems to be giving an ip using dhcp. But after some time, bridge ip is changing to below address */
	bridge, err := Dial("192.168.0.100:8899")
	//fmt.Println(bridge)
	if err != nil {
		fmt.Println("Something wrong")
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("Provide atleast one argument")
		return
	}
	option := os.Args[1]

	switch option {
	case "all_off":
		bridge.SendCommand(ALL_OFF)
	case "all_on":
		bridge.SendCommand(ALL_ON)
	case "all_disco":
		bridge.SendCommand(ALL_DISCO)
	case "all_white":
		bridge.SendCommand(ALL_WHITE)
		bridge.SendCommand(ALL_ON_FULL)
	default:
		fmt.Println("Wrong option. Provide some thing appropriate")
	}
}
