package backend

import (
	"fmt"
	"time"
)

type Setting struct {
	SettingKey   string
	SettingValue string
}

var settings []Setting

func ReceiveSettings(keys []string, values []string) {
	settings = nil
	for i := 0; i < len(keys); i++ {
		var newElem Setting
		newElem.SettingKey = keys[i]
		newElem.SettingValue = values[i]
		settings = append(settings, newElem)
	}
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func printStatusToTerminal(t time.Time) {
	fmt.Println("[daemon tick]")
	fmt.Printf("%+v\n", settings)
	fmt.Println("[daemon tock]")
}

func TransmissionDaemon() {
	doEvery(10*time.Second, printStatusToTerminal)
}
