package main

import (
	"fmt"
	"log"

	"github.com/ninjasphere/go-ninja/support"
)

func main() {
	fmt.Println("RTR. NinjaSphere LimitlessLed driver")
	_, err := NewLimitlessLedDriver()
	if err != nil {
		log.Fatalf("Failed to create fake driver: %s", err)
	}

	support.WaitUntilSignal()
}
