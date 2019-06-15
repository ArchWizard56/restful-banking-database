package main

import (
	"fmt"
	"os"
	//	"github.com/gorilla/mux"
	"encoding/json"
	"flag"
	"log"
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
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config
}

func main() {
	fmt.Println("Loading Config")
    config := loadConfig()
	fmt.Println("Initializing Router")
    router := InitRouter()
    fmt.Println("Serving http")
    err := http.ListenAndServe(":"+strconv.Itoa(config.Port), router)
    if err != nil {
        log.Fatal(err)
    }
}
