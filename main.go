package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-ping/ping"
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

	// Checks if any targets are configured
	if len(Config.PingTargets) == 0 {
		checkError(errors.New("no ping targets configured"))
	}

	// Iterates over the PingTargets and spawns a pinger instance for each
	fmt.Println("Creating pinger instances")
	for i := range Config.PingTargets {
		go pinger(&Config.PingTargets[i])
	}

	// Sleeps to keep application open
	time.Sleep(20 * time.Second)
}

// Takes a pingTarget, checks the validity then pings the target
// Currently just dumps the data to stdout
// Meant to be used as a goroutine
func pinger(target *pingTarget) {
	// Checks if host is set, otherwise quits
	if target.Host == "" {
		checkError(errors.New("one of the ping targets is missing a hostname/ip"))
	}

	// Checks if name is set, defaults to hostname if its not
	if target.Name == "" {
		fmt.Println("Missing name for one of ping targets, defaulting to hostname")
		target.Name = target.Host
	}

	// Default timeout setting
	if target.Timeout == 0 {
		target.Timeout = 3
	}

	// Default colour setting
	if target.Colour == "" {
		target.Colour = "purple"
	}

	// Default period setting
	if target.Period == 0 {
		target.Period = 5
	}

	timeout := time.Duration(target.Timeout) * time.Second

	// Pings forever
	for {
		// Have to recreate pinger?
		Pinger, err := ping.NewPinger(target.Host)
		checkError(err)

		Pinger.Timeout = timeout
		err = Pinger.Run()
		checkError(err)

		fmt.Println(target.Host, Pinger.Statistics())
		time.Sleep(time.Second * time.Duration(target.Period))
	}
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
