package ServiceRegistry

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/hariadivicky/nano" // import RESTful methods.
	_ "github.com/mattn/go-sqlite3"
)

func StartServiceRegistry() {

	config := getConfig()

	OpenDatabase("file:ServiceRegistry/database/registryDB.db?cache=private&mode=rwc&_foreign_keys=on")
	if config.CleanEndOfValidity {
		go startValidityTimer(config.CleanDelay) //starts cleaning on ticks in background
	}
	n := nano.New()

	n.Use(nano.Recovery())

	n.Use(nano.CORSWithConfig(nano.CORSConfig{AllowedOrigins: []string{"*"}}))
	method := new(JsonFile)

	n.GET("/serviceregistry/echo", method.Echo)
	n.POST("/serviceregistry/query", method.Query)
	n.POST("/serviceregistry/register", method.Store)
	n.DELETE("/serviceregistry/unregister", method.Unregister)
	log.Println("server running at port: " + strconv.Itoa(config.Port))

	err := n.Run(":" + strconv.Itoa(config.Port))
	if err != nil {
		log.Fatalf("could not start application: %v", err)
	}
}

type configJson struct {
	Port               int  `json:"port"`
	CleanEndOfValidity bool `json:"cleanEndOfValidity"`
	CleanDelay         int  `json:"endofValidityCleaningCycle"`
}

func getConfig() configJson {
	dat, err := ioutil.ReadFile("./config.json")
	var config configJson
	if err != nil {
		println("No config found. Using default parameters.")
		config.Port = 4245
		config.CleanEndOfValidity = true
		config.CleanDelay = 10

		dat, err := json.Marshal(config)
		if err != nil {
			panic(err.Error())
		}
		err = ioutil.WriteFile("./config.json", dat, 0644)
		if err != nil {
			panic(err.Error())
		}

	} else {
		err = json.Unmarshal(dat, &config)
		if err != nil {
			panic(err.Error())
		}

	}

	return config
}
