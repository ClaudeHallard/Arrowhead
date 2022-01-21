package main

import (
	"log"

	"github.com/hariadivicky/nano" // import RESTful methods.
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//OpenDatabase("file:registryDB.db?cache=private")
	OpenDatabase("file:registryDB.db?cache=private&_foreign_keys=on")

	n := nano.New()

	n.Use(nano.Recovery())

	n.Use(nano.CORSWithConfig(nano.CORSConfig{AllowedOrigins: []string{"*"}}))

	method := new(JsonFile)
	n.GET("/serviceregistry/echo", method.Echo)
	n.POST("/serviceregistry/query", method.Query)
	n.POST("/serviceregistry/register", method.Store)
	n.DELETE("/serviceregistry/unregister", method.Unregister)

	log.Println("server running at port 4245")

	err := n.Run(":4245")
	if err != nil {
		log.Fatalf("could not start application: %v", err)
	}
}
