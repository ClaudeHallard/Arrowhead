package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ServiceRegistryEntry struct {
	ServiceDefinition string `json:"serviceDefinition"`
	ProviderSystem    ProviderSystem
	ServiceURI        string `json:"serviceURI"`
	EndOfvalidity     string `json:"endOfValidity"`
	Secure            string `json:"NOT_SECURE"`
	Metadata          Metadata
	Version           int      `json:"version"`
	Interfaces        []string `json:"interfaces"`
}
type ProviderSystem struct {
	SystemName         string `json:"systemName"`
	Address            string `json:"adress"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
}

type Metadata struct {
	AdditionalProp1 string `json:"additionalProp1"`
	AdditionalProp2 string `json:"additionalProp2"`
	AdditionalProp3 string `json:"additionalProp3"`
}

func main() {
	serviceRegistryEntry := &ServiceRegistryEntry{
		ServiceDefinition: "aa",
		ProviderSystem: ProviderSystem{
			SystemName:         "bb",
			Address:            "cc",
			Port:               222,
			AuthenticationInfo: "dd",
		},
		ServiceURI:    "ee",
		EndOfvalidity: "ff",
		Secure:        "gg",
		Metadata: Metadata{
			AdditionalProp1: "hh",
			AdditionalProp2: "ii",
			AdditionalProp3: "jj",
		},

		Version: 33,
		Interfaces: []string{
			"Interface1",
			"Interface2",
			"Interface3",
			"Interface4",
		},
	}

	body, _ := json.Marshal(serviceRegistryEntry)
	resp, err := http.Post("http://127.0.0.1:4245/serviceregistry/register", "application/json", bytes.NewBuffer(body))
	if err != nil {

		panic(err)
	}
	defer resp.Body.Close()

	//Check response code, (201)
	if resp.StatusCode == http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Response: ")
			//Failed to read response.
			panic(err)
		}

		//Convert bytes to String and print
		jsonStr := string(body)
		fmt.Println("Response: ", jsonStr)

	} else {
		//The status is not Created. print the error.
		fmt.Println("Get failed with error: ", resp.Status)
	}
}
