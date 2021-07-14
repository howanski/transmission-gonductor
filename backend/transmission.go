package backend

import (
	"fmt"
	"time"
)

// type Configs struct {
// 	ConfigKey   string
// 	ConfigValue string
// }

// var configs []ConfigStorage

// func refreshConfigs() {
// 	//TODO :)
// }

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func printStatusToTerminal(t time.Time) {
	fmt.Printf("%v: Daemon active..\n", t)
}

func TransmissionDaemon() {
	doEvery(10*time.Second, printStatusToTerminal)
}
