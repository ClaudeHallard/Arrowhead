package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	values := map[string]string{
		"serviceDefinition":  "aa",
		"systemName":         "bb",
		"metaData":           "metaData",
		"port":               "Port",
		"authenticationInfo": "authenticationInfo",
		"serviceURI":         "serviceURI",
		"endOfValidity":      "endOfValidity",
		"secure":             "secure",
		"address":            "address",
		"version":            "version",
		"interfaces":         "interfaces"}
	println("Sending...")
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://127.0.0.1:4245/serviceregistry/register", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res["json"])
}
