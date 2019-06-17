package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
    "database/sql"
)

type Configuration struct {
	Port int
    Database string
}
type Secrets struct {
    Jwtsecret string
}

var configFileLocation *string
var secretFileLocation *string
var debug *bool
var Database *sql.DB
var Secret Secrets
var Config Configuration

func SetFlags() {
	configFileLocation = flag.String("c", "config.json", "the location of the config file to use")
	secretFileLocation = flag.String("s", "secrets.json", "the location of the secrets file to use")
	debug = flag.Bool("d", false, "toggles debug output")
	flag.Parse()
}

func loadConfig () (Configuration, Secrets) {
	DualInfo("Loading Config")
	var config Configuration
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
    file, err = os.Open(*secretFileLocation)
	if err != nil {
		DualErr(err)
	}
    var secrets Secrets
	//Decode the json config file
	decoder = json.NewDecoder(file)
	err = decoder.Decode(&secrets)
	if err != nil {
		DualErr(err)
	}
	return config, secrets
}

func main() {
	SetFlags()
	SetupLogging(*debug)
	DualInfo("Initialized Logging")
	Config, Secret = loadConfig()
	DualInfo("Loading Database")
    Database = LoadDB(Config.Database)
	DualInfo("Initializing Router")
	router := InitRouter()
	DualInfo("Starting Server")
	//Use the router we have to listen and serve http
	err := http.ListenAndServe(":"+strconv.Itoa(Config.Port), router)
	if err != nil {
		DualErr(err)
	}
}
