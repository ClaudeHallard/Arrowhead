package main

import (
	"log"

	"github.com/hariadivicky/nano" // import RESTful methods.
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	OpenDatabase("file:registryDB.db?cache=private")

	n := nano.New()

	n.Use(nano.Recovery())

	n.Use(nano.CORSWithConfig(nano.CORSConfig{AllowedOrigins: []string{"*"}}))

	method := new(JsonFile)
	n.GET("/serviceregistry/echo", method.Echo)
	n.POST("/serviceregistry/query", method.Query)
	n.POST("/serviceregistry/register", method.Store)
	n.DELETE("/serviceregistry/unregister", method.Unregister)

	log.Println("server running at port 8080")

	err := n.Run(":8080")
	if err != nil {
		log.Fatalf("could not start application: %v", err)
	}
}
