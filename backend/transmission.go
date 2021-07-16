package backend

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Statistic struct {
	StatisticKey   string
	StatisticValue string
}

var statistics []Statistic

type Setting struct {
	SettingKey   string
	SettingValue string
}

var settings []Setting

var transmissionCrossKey string

var spamTerminal bool

func getHumanReadableTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func ReceiveSettings(keys []string, values []string) {
	settings = nil
	for i := 0; i < len(keys); i++ {
		var newElem Setting
		newElem.SettingKey = keys[i]
		newElem.SettingValue = values[i]
		settings = append(settings, newElem)
	}
	if findInSettings("gonductorDebugToTerminal") == "yes" {
		spamTerminal = true
	} else {
		spamTerminal = false
	}
}

func GiveStatistics() map[string]interface{} {
	data := make(map[string]interface{})
	for i := 0; i < len(statistics); i++ {
		data[statistics[i].StatisticKey] = statistics[i].StatisticValue
	}
	return data
}

func findInSettings(key string) string {
	for i := 0; i < len(settings); i++ {
		elem := settings[i]
		if elem.SettingKey == key {
			return elem.SettingValue
		}
	}
	return ""
}

func saveStatistic(statisticKey string, value string) {
	statFound := false
	for i := 0; i < len(statistics); i++ {
		if statistics[i].StatisticKey == statisticKey {
			statFound = true
			statistics[i].StatisticValue = value
		}
	}
	if !statFound {
		newStat := Statistic{StatisticKey: statisticKey, StatisticValue: value}
		statistics = append(statistics, newStat)
	}
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func makeTransmissionRequest(json_data []byte) []byte {
	if spamTerminal {
		fmt.Println("----------------------------------------> SENDING: " + string(json_data))
	}

	userName := findInSettings("transmissionUser")
	password := findInSettings("transmissionPassword")
	host := findInSettings("transmissionHost")
	var fullReqString = "http://" + host + ":9091/transmission/rpc/"

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest("POST", fullReqString, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Printf("[1] Got error %s", err.Error())
	} else {
		req.SetBasicAuth(userName, password)
		req.Header.Set("X-Transmission-Session-Id", transmissionCrossKey)
		response, err := client.Do(req)
		if err != nil {
			saveStatistic("connectionStatus", "Transmission server GOne")
			if spamTerminal {
				fmt.Printf("[No response] Got error %s", err.Error())
			}

		} else {
			if response.StatusCode == 409 {
				saveStatistic("connectionStatus", "Updating CORS")
				saveStatistic("lastPing", getHumanReadableTime())
				transmissionCrossKey = response.Header.Get("X-Transmission-Session-Id")
				if spamTerminal {
					fmt.Println("Cross-dressing for CSRF: " + transmissionCrossKey)
				}
				deepBody := makeTransmissionRequest(json_data)
				return deepBody
			} else if response.StatusCode == 401 {
				saveStatistic("connectionStatus", "Authorization error")
				saveStatistic("lastPing", getHumanReadableTime())
				if spamTerminal {
					fmt.Println("--------------- [ TRANSMISSION AUTHORIZATION ERROR ] ---------------")
				}
			} else {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					fmt.Printf("[3] Got error %s", err.Error())
				} else {
					saveStatistic("connectionStatus", "Operational")
					saveStatistic("lastPing", getHumanReadableTime())
					sb := string(body)
					if spamTerminal {
						fmt.Print("RESPONSE: ------------------------> " + sb)
					}
					defer response.Body.Close()
					return body
				}
			}
		}
	}
	var failResponse []byte
	return failResponse
}

func getGeneralTorrentsData() {
	generalQuestion := []byte(`{"method":"torrent-get","arguments":{"fields": ["id", "name", "totalSize"]}}`)
	makeTransmissionRequest(generalQuestion)
}

func makeOrchestratedRound(t time.Time) {
	fmt.Println("[GONDUCTOR PULSE] [" + getHumanReadableTime() + "]")

	if spamTerminal {
		fmt.Println("[daemon tick] starting functions round")
	}
	getGeneralTorrentsData()
	if spamTerminal {
		fmt.Println("[daemon tick] round end")
	}
}

func TransmissionDaemon() {
	doEvery(10*time.Second, makeOrchestratedRound)
}
