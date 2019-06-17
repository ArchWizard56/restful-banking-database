package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
    "database/sql"
)

type Config struct {
	Port int
    Database string
}

var configFileLocation *string
var debug *bool
var Database *sql.DB

func SetFlags() {
	configFileLocation = flag.String("c", "config.json", "the location of the config file to use")
	debug = flag.Bool("d", false, "toggles debug output")
	flag.Parse()
}

func loadConfig() Config {
	DualInfo("Loading Config")
	var config Config
	//Load the location of the config from the command line, with the default of "config.json"
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
	SetFlags()
	SetupLogging(*debug)
	DualInfo("Initialized Logging")
	config := loadConfig()
	DualInfo("Loading Database")
    Database = LoadDB(config.Database)
	DualInfo("Initializing Router")
	router := InitRouter()
	DualInfo("Starting Server")
	//Use the router we have to listen and serve http
	err := http.ListenAndServe(":"+strconv.Itoa(config.Port), router)
	if err != nil {
		DualErr(err)
	}
}
