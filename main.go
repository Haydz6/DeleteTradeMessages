package main

import (
	"encoding/json"
	"os"
	"path"
)

type SettingsStruct struct {
	DeleteTradeReceived  bool
	DeleteTradeAccepted  bool
	DeleteTradeDeclined  bool
	DeleteTradeCountered bool
}

var ROBLOSECURITY string
var Settings SettingsStruct

func FetchROBLOSecurity() {
	Bytes, err := os.ReadFile(path.Join(".", "ROBLOSECURITY.txt"))

	if err != nil {
		panic(err.Error())
	}

	ROBLOSECURITY = string(Bytes)
}

func FetchSettings() {
	Bytes, err := os.ReadFile(path.Join(".", "settings.json"))

	if err != nil {
		panic(err.Error())
	}

	jsonErr := json.Unmarshal(Bytes, &Settings)

	if jsonErr != nil {
		panic(jsonErr.Error())
	}
}

func main() {
	FetchROBLOSecurity()
	FetchSettings()

	ScanMessages()
}
