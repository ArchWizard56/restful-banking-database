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
    //Load the location of the config from the command line, with the default of "config.json"
	configFileLocation := flag.String("c", "config.json", "the location of the config file to use")
    //Try and open it, if we can't, exit with a fatal error
	file, err := os.Open(*configFileLocation)
	if err != nil {
		DualErr(err)
	}
    //Decode the json config file
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
    //Use the router we have to listen and serve http
    err := http.ListenAndServe(":"+strconv.Itoa(config.Port), router)
    if err != nil {
        DualErr(err)
    }
}
