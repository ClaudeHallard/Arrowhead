package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	//echoTest() // Test example to send echo call
	//registerTest() // Test example to register a service
	//queryTest()
	//unregisterTest()

}

//createAndRegister(1)
//createAndRegister(2)
//createAndRegister(3)

//Shared Structs
//Metadata scruct, used by
//ServiceRegistryEntryInput, ServiceRegistryEntryOutput,
// ServiceQueryForm and ServiceQueryList.
type ServiceDefinition struct {
	ID                int    `json:"id"`
	ServiceDefinition string `json:"serviceDefinition"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}
type Provider struct {
	ID                 int    `json:"id"`
	SystemName         string `json:"systemName"`
	Address            string `json:"adress"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type Interface struct {
	ID            int    `json:"id"`
	InterfaceName string `json:"interfaceName"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// ServiceRegistryEntry Input Version
type ServiceRegistryEntryInput struct {
	ServiceDefinition string         `json:"serviceDefinition"`
	ProviderSystem    ProviderSystem `json:"providerSystem"`
	ServiceUri        string         `json:"serviceUri"`
	EndOfvalidity     string         `json:"endOfValidity"`
	Secure            string         `json:"NOT_SECURE"`
	Metadata          []string       `json:"metadata"`
	Version           int            `json:"version"`
	Interfaces        []string       `json:"interfaces"`
}
type ProviderSystem struct {
	SystemName         string `json:"systemName"`
	Address            string `json:"adress"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
}

type ServiceUnregisterEntryInput struct {
	ServiceDefinition string `json:"serviceDefinition"`
	SystemName        string `json:"systemName"`
	Address           string `json:"address"`
	Port              int    `json:"port"`
}

// ServiceRegistryEntry Output Version
type ServiceRegistryEntryOutput struct {
	ID                int               `json:"id"`
	ServiceDefinition ServiceDefinition `json:"serviceDefinition"`
	Provider          Provider          `json:"provider"`
	ServiceUri        string            `json:"serviceUri"`
	EndOfValidity     string            `json:"endOfValidity"`
	Secure            string            `json:"NOT_SECURE"`
	Metadata          []string          `json:"metadata"`
	Version           int               `json:"version"`
	Interfaces        []Interface       `json:"interfaces"` //I think this is how it's implemented -Ivar (should result in an array of interfaces.)
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
}

//ServiceQueryForm
type ServiceQueryForm struct {
	ServiceDefinitionRequirement string   `json:"serviceDefinitionRequirement"`
	InterfaceRequirements        []string `json:"interfaceRequirements"`
	SecurityReRequirements       []string `json:"securityRequirements"`
	MetadataRequirements         []string `json:"metadataRequirements "`
	VersionRequirements          int      `json:"versionRequirements"`
	MaxVersionRequirements       int      `json:"maxVersionRequirements"`
	MinVersionRRequirements      int      `json:"minVersionRRequirements"`
	PingProviders                bool     `json:"pingProviders"`
}

//ServiceQueryList
type ServiceQueryList struct {
	ServiceQueryData []ServiceRegistryEntryOutput `json:"serviceQueryData"`
	UnfilteredHits   int                          `json:"unfilteredHits"`
}

// Package Handler for the POST method
func sendPackage(body []byte, path string) {
	resp, err := http.Post("http://127.0.0.1:4245/serviceregistry/"+path, "application/json", bytes.NewBuffer(body))
	if err != nil {

		panic(err)
	}
	defer resp.Body.Close()

	//Check response code, (201)

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Response: ")
		//Failed to read response.
		panic(err)
	}

	//The status is not Created. print the error.
	fmt.Println("	Respons status: ", resp.Status)
	//Convert bytes to String and print
	jsonStr := string(body)
	fmt.Println("	Response: ", jsonStr)

}

// example function to test the echo request
func echoTest() {
	resp, err := http.Get("http://127.0.0.1:4245/serviceregistry/echo")
	if err != nil {
		println("Echo did not work")
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	println(string(body))
}

// example structure on how to register a service into the serviceRegistry
func registerTest() {

	serviceRegistryEntry := &ServiceRegistryEntryInput{
		ServiceDefinition: "Temperature",
		ProviderSystem: ProviderSystem{
			SystemName:         "TempProviderLTU",
			Address:            "127.0.0.1",
			Port:               9999,
			AuthenticationInfo: "NOT_USED",
		},
		ServiceUri:    "http://127.0.0.1:9999/TempProviderLTU/Temperature",
		EndOfvalidity: "2022-01-25 15:04:05",
		Secure:        "NOT_SECURE", // tokens not used
		Metadata: []string{
			"Celsius",
			"Kelvin",
			"Fahrenheit",
		},

		Version: 1,
		Interfaces: []string{
			"JSON",
			"HTTP-SECURE-JSON",
		},
	}

	body, _ := json.Marshal(serviceRegistryEntry)
	sendPackage(body, "register")

}

func queryTest() {
	serviceQueryForm := ServiceQueryForm{
		ServiceDefinitionRequirement: "Temperature",
		InterfaceRequirements:        []string{},
		SecurityReRequirements:       []string{},
		MetadataRequirements:         []string{},
		VersionRequirements:          0,
		MaxVersionRequirements:       0,
		MinVersionRRequirements:      0,
		PingProviders:                false,
	}
	body, _ := json.Marshal(serviceQueryForm)
	sendPackage(body, "query")

}

func unregisterTest() {

	println("test")
}

// func createAndRegister(i int) {
// 	serviceRegistryEntry := &ServiceRegistryEntryInput{
// 		ServiceDefinition: "serviceDef" + strconv.Itoa(i),
// 		ProviderSystem: ProviderSystem{
// 			SystemName:         "systemName" + strconv.Itoa(i),
// 			Address:            "address" + strconv.Itoa(i),
// 			Port:               i,
// 			AuthenticationInfo: "authInfo" + strconv.Itoa(i),
// 		},
// 		ServiceUri:    "serviceUri" + strconv.Itoa(i),
// 		EndOfvalidity: "endofValidity" + strconv.Itoa(i),
// 		Secure:        "NOT_SECURE",
// 		Metadata: []string{
// 			strconv.Itoa(i) + "metadataA",
// 			strconv.Itoa(i) + "metadataB",
// 			strconv.Itoa(i) + "metadataC",
// 			strconv.Itoa(i) + "metadataD",
// 		},

// 		Version: i,
// 		Interfaces: []string{
// 			strconv.Itoa(i) + "InterfaceA",
// 			strconv.Itoa(i) + "InterfaceB",
// 			strconv.Itoa(i) + "InterfaceC",
// 			strconv.Itoa(i) + "InterfaceD",
// 		},
// 	}
// 	body, _ := json.Marshal(serviceRegistryEntry)
// 	sendPackage(body, "register")

// }
