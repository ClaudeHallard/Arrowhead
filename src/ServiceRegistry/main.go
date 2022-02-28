/*
**********************************************
@authors:
Ivar Grunaeu, ivagru-9@student.ltu.se
Jean-Claude, Hallard jeahal-8@student.ltu.se
Marcus Paulsson, marpau-8@student.ltu.se
Pontus Schünemann, ponsch-9@student.ltu.se
Fabian Widell, fabwid-9@student.ltu.se

**********************************************
main.go contains the main functions for starting the service registry, create settings
for the server and a listener to a port. */

package ServiceRegistry

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/hariadivicky/nano"  // Import RESTful methods.
	_ "github.com/mattn/go-sqlite3" // Import SQLite methods.
)

// Boot server and starts to listen for http requests.
func StartServiceRegistry() {

	config := getConfig()

	OpenDatabase()
	if config.CleanEndOfValidity {
		go startValidityTimer(config.CleanDelay) // Starts cleaning on ticks in background as a go routine.
	}

	n := nano.New()

	n.Use(nano.Recovery()) // Recovers the server if errors occurs

	// Supported requests:
	n.GET("/serviceregistry/echo", Echo)
	n.POST("/serviceregistry/query", Query)
	n.POST("/serviceregistry/register", Store)
	n.DELETE("/serviceregistry/unregister", Unregister)

	// Set the program to listen to the specific port
	log.Println("server running at port: " + strconv.Itoa(config.Port))
	err := n.Run(":" + strconv.Itoa(config.Port))
	if err != nil {
		log.Fatalf("could not start application: %v", err)
	}

}

// Go struct corresponding to current config file.
type configJson struct {
	Port               int  `json:"port"`
	CleanEndOfValidity bool `json:"cleanEndOfValidity"`
	CleanDelay         int  `json:"endofValidityCleaningCycle"`
}

// Returns the current data read from the config file
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
