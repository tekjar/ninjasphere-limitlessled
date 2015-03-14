package main

import (
	"fmt"
	"github.com/ninjasphere/go-ninja/support"
	"os"
	"os/signal"
)

func main() {
	NewLimitlessLedDriver()
	support.WaitUntilSignal()
}
