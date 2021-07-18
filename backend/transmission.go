package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

type TorrentInsideFile struct {
	FileName       string
	FileLength     int
	FileLengthDone int
	FileIsWanted   bool
	FileId         string
}

type Torrent struct {
	TorrentId    string
	TorrentName  string
	TorrentFiles []TorrentInsideFile
}

var torrentStorage []Torrent

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
	hostPort := findInSettings("transmissionHostPort")
	var fullReqString = "http://" + host + ":" + hostPort + "/transmission/rpc/"

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

func gutterJsonInterfaceStringIndexed(outerInterface interface{}, key string) interface{} {
	mapmap := outerInterface.(map[string]interface{})
	innerInterface := mapmap[key]
	return innerInterface
}

func gutterJsonInterfaceIntegerIndexed(outerInterface interface{}, key int) interface{} {
	mapmap := outerInterface.([]interface{})
	innerInterface := mapmap[key]
	return innerInterface
}

func jsonBytesToInterface(bodyBytes []byte) interface{} {
	var interfaceWithJson interface{}
	json.Unmarshal(bodyBytes, &interfaceWithJson)
	return interfaceWithJson
}

func getLengthOfArrayInInterface(face interface{}) int {
	faceArr := face.([]interface{})
	return len(faceArr)
}

func faceToString(x interface{}) string {
	str := fmt.Sprintf("%v", x)
	return str
}

func faceToFloat(face interface{}) float64 {
	float := face.(float64)
	return float
}

func faceToInt(face interface{}) int {
	float := faceToFloat(face)
	inter := int(float)
	return inter
}

func getGeneralTorrentsData() bool {
	generalQuestion := []byte(`{"method":"torrent-get","arguments":{"fields": ["id", "name", "totalSize", "sizeWhenDone", "files", "fileStats", "rateDownload", "isFinished", "isStalled", "eta", "priorities", "wanted"]}}`)
	jsonInterfaceResponse := jsonBytesToInterface(makeTransmissionRequest(generalQuestion))
	if jsonInterfaceResponse != nil {
		arguments := gutterJsonInterfaceStringIndexed(jsonInterfaceResponse, "arguments")
		torrents := gutterJsonInterfaceStringIndexed(arguments, "torrents")
		torrentStorage = nil
		for i := 0; i < getLengthOfArrayInInterface(torrents); i++ {
			torrent := gutterJsonInterfaceIntegerIndexed(torrents, i)
			files := gutterJsonInterfaceStringIndexed(torrent, "files")
			fileStats := gutterJsonInterfaceStringIndexed(torrent, "fileStats")
			fileStatsLength := getLengthOfArrayInInterface(fileStats)
			if fileStatsLength > 0 { // Any files known
				var torrentFilez []TorrentInsideFile
				for j := 0; j < fileStatsLength; j++ {
					fileData := gutterJsonInterfaceIntegerIndexed(files, j)
					fileStatData := gutterJsonInterfaceIntegerIndexed(fileStats, j)
					isWanted := faceToString(gutterJsonInterfaceStringIndexed(fileStatData, "wanted"))
					var terrentFile TorrentInsideFile
					terrentFile.FileName = faceToString(gutterJsonInterfaceStringIndexed(fileData, "name"))
					terrentFile.FileLength = faceToInt(gutterJsonInterfaceStringIndexed(fileData, "length"))
					terrentFile.FileLengthDone = faceToInt(gutterJsonInterfaceStringIndexed(fileData, "bytesCompleted"))
					terrentFile.FileIsWanted = isWanted == "true"
					terrentFile.FileId = fmt.Sprint(j)
					torrentFilez = append(torrentFilez, terrentFile)
				}

				var storageRecord = Torrent{}
				storageRecord.TorrentName = faceToString(gutterJsonInterfaceStringIndexed(torrent, "name"))
				storageRecord.TorrentId = faceToString(gutterJsonInterfaceStringIndexed(torrent, "id"))
				storageRecord.TorrentFiles = torrentFilez
				torrentStorage = append(torrentStorage, storageRecord)
			}
		}
		return true
	}
	return false
}

func updatePrioritiesOnSubFiles() {
	ifManage := findInSettings("transmissionManagePrioritiesAlphabetically")
	if ifManage == "yes" {
		for i := 0; i < len(torrentStorage); i++ {
			torrentDescription := torrentStorage[i]
			torrentFiles := torrentDescription.TorrentFiles
			var foundHighPriority bool = false
			var highOnes []TorrentInsideFile
			var lowOnes []TorrentInsideFile
			sort.Slice(torrentFiles[:], func(i, j int) bool {
				return strings.ToLower(torrentFiles[i].FileName) < strings.ToLower(torrentFiles[j].FileName)
			})
			for j := 0; j < len(torrentFiles); j++ {
				torrentFile := torrentFiles[j]
				torrentFilePriorityShouldBeHigh := torrentFile.FileIsWanted && (torrentFile.FileLength != torrentFile.FileLengthDone)
				if torrentFilePriorityShouldBeHigh {
					if foundHighPriority {
						lowOnes = append(lowOnes, torrentFile)
					} else {
						highOnes = append(highOnes, torrentFile)
						foundHighPriority = true
					}
				} else {
					lowOnes = append(lowOnes, torrentFile)
				}
			}
			var lowIds string
			for k := 0; k < len(lowOnes); k++ {
				lowFile := lowOnes[k]
				lowIds += lowFile.FileId
				if k < len(lowOnes)-1 {
					lowIds += ","
				}
			}
			var highIds string
			for k := 0; k < len(highOnes); k++ {
				highFile := highOnes[k]
				highIds += highFile.FileId
				if k < len(highOnes)-1 {
					highIds += ","
				}

			}
			instruction := []byte(`{"method":"torrent-set","arguments":{"ids":[` + torrentDescription.TorrentId + `],"priority-low":[` + lowIds + `],"priority-high":[` + highIds + `]}}`)
			makeTransmissionRequest(instruction)

		}
	}
}

func makeOrchestratedRound(t time.Time) {
	fmt.Println("[GONDUCTOR PULSE] [" + getHumanReadableTime() + "]")

	if spamTerminal {
		fmt.Println("[daemon tick] starting functions round")
	}
	success := getGeneralTorrentsData()
	if success {
		updatePrioritiesOnSubFiles()
	}
	if spamTerminal {
		fmt.Println("[daemon tick] round end")
	}
}

func TransmissionDaemon() {
	doEvery(10*time.Second, makeOrchestratedRound)
}
