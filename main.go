package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type pingTarget struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	Timeout int    `json:"timeout"`
	Colour  string `json:"colour"`
	Period  int    `json:"period"`
}

type config struct {
	ServerIP    string       `json:"serverIP"`
	ServerPort  int          `json:"serverPort"`
	PingTargets []pingTarget `json:"pingTargets"`
}

func main() {
	// Load configuration file - conf.json
	var Config config
	loadConfig(&Config)
	fmt.Println("Loaded config")

	fmt.Println(Config.PingTargets)
}

// Function to load configuration file into config struct
// 	- Doesn't add default values
// 	- conf.json only for now
func loadConfig(Config *config) {
	conf, err := os.Open("conf.json")
	if err != nil {
		fmt.Println("Failed to open configuration file")
		checkError(err)
	}
	defer conf.Close()

	configBytes, err := ioutil.ReadAll(conf)
	checkError(err)

	json.Unmarshal(configBytes, &Config)
}

// Quick helper function to save my fingers
// When err != nil, dumps error and panics
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
