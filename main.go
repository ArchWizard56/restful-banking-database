package main

import (
	"os"
	"encoding/json"
	"flag"
    "net/http"
    "strconv"
)

type Config struct {
	Port int
}

func loadConfig() Config {
    DualInfo("Loading Config")
	var config Config
	configFileLocation := flag.String("c", "config.json", "the location of the config file to use")
	file, err := os.Open(*configFileLocation)
	if err != nil {
		DualErr(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		DualErr(err)
	}
	return config
}

func main() {
    SetupLogging()
    DualInfo("Initialized Logging")
    config := loadConfig()
	DualInfo("Initializing Router")
    router := InitRouter()
    DualInfo("Starting Server")
    err := http.ListenAndServe(":"+strconv.Itoa(config.Port), router)
    if err != nil {
        DualErr(err)
    }
}
