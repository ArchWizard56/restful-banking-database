package main

import (
	//"fmt"
	"os"
	//	"github.com/gorilla/mux"
	"encoding/json"
	"flag"
    "net/http"
    "strconv"
)

type Config struct {
	Port int
}

func loadConfig() Config {
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
    DualInfo("Loading Config")
    config := loadConfig()
	DualInfo("Initializing Router")
    router := InitRouter()
    DualInfo("Serving http")
    err := http.ListenAndServe(":"+strconv.Itoa(config.Port), router)
    if err != nil {
        DualErr(err)
    }
}
