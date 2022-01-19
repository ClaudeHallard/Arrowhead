package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	//registerTest()
	//createAndRegister(1)
	//createAndRegister(2)
	//createAndRegister(3)
	queryTest()

}

//Shared Structs

//Metadata scruct, used by
//ServiceRegistryEntryInput, ServiceRegistryEntryOutput,
// ServiceQueryForm and ServiceQueryList.
type MetadataOld struct {
	AdditionalProp1 string `json:"additionalProp1"`
	AdditionalProp2 string `json:"additionalProp2"`
	AdditionalProp3 string `json:"additionalProp3"`
}

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

func registerTest() {
	serviceRegistryEntry := &ServiceRegistryEntryInput{
		ServiceDefinition: "aa",
		ProviderSystem: ProviderSystem{
			SystemName:         "bb",
			Address:            "cc",
			Port:               222,
			AuthenticationInfo: "dd",
		},
		ServiceUri:    "ee",
		EndOfvalidity: "ff",
		Secure:        "gg",
		Metadata: []string{
			"metadata1",
			"metadata2",
			"metadata3",
			"metadata4",
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
	sendPackage(body, "register")

}

func createAndRegister(i int) {

	serviceRegistryEntry := &ServiceRegistryEntryInput{
		ServiceDefinition: "serviceDef" + strconv.Itoa(i),
		ProviderSystem: ProviderSystem{
			SystemName:         "systemName" + strconv.Itoa(i),
			Address:            "address" + strconv.Itoa(i),
			Port:               i,
			AuthenticationInfo: "authInfo" + strconv.Itoa(i),
		},
		ServiceUri:    "serviceUri" + strconv.Itoa(i),
		EndOfvalidity: "endofValidity" + strconv.Itoa(i),
		Secure:        "NOT_SECURE",
		Metadata: []string{
			strconv.Itoa(i) + "metadataA",
			strconv.Itoa(i) + "metadataB",
			strconv.Itoa(i) + "metadataC",
			strconv.Itoa(i) + "metadataD",
		},

		Version: i,
		Interfaces: []string{
			strconv.Itoa(i) + "InterfaceA",
			strconv.Itoa(i) + "InterfaceB",
			strconv.Itoa(i) + "InterfaceC",
			strconv.Itoa(i) + "InterfaceD",
		},
	}
	body, _ := json.Marshal(serviceRegistryEntry)
	sendPackage(body, "register")

}
func queryTest() {
	serviceQueryForm := ServiceQueryForm{
		ServiceDefinitionRequirement: "serviceDef1",
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
