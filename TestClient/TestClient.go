package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	convert = false //Set to true to convert from metadata array to struct
	println("Converting when sending: " + strconv.FormatBool(convert))

	//preformTests("query", 10000, 40) //preforms 100 querys
	preformTests("register", 5, 100)
}

var convert bool //Set to true to convert from metadata array to struct

func (metadata Metadata) MarshalJSON() ([]byte, error) {
	metadataStruct := MetadataStruct{}
	if convert {
		if len(metadata) >= 1 {
			metadataStruct.AdditionalProp1 = metadata[0]
		}
		if len(metadata) >= 2 {
			metadataStruct.AdditionalProp2 = metadata[1]
		}
		if len(metadata) >= 3 {
			metadataStruct.AdditionalProp3 = metadata[2]
		}

		return json.Marshal(metadataStruct)
	} else {
		var sArray []string
		sArray = metadata
		return json.Marshal(sArray)
	}

}
func (metadata *Metadata) UnmarshalJSON(data []byte) error {

	var metadataStruct *MetadataStruct
	println("outsidesss22")
	var stringArray []string
	println(len(data))
	println(string(data))
	err := json.Unmarshal(data, &stringArray)

	if err != nil {
		println("inside")
		err := json.Unmarshal(data, &metadataStruct)
		if err != nil {
			return err
		}
		println(metadataStruct.AdditionalProp1)

		if metadataStruct.AdditionalProp1 != "" {
			*metadata = append(*metadata, metadataStruct.AdditionalProp1)
		}
		if metadataStruct.AdditionalProp2 != "" {
			*metadata = append(*metadata, metadataStruct.AdditionalProp2)
		}
		if metadataStruct.AdditionalProp3 != "" {
			*metadata = append(*metadata, metadataStruct.AdditionalProp3)
		}

		return nil

	}
	*metadata = append(stringArray)
	println("outside")
	//*metadata = append(stringArray)
	return nil
}

type testData struct {
	status      string
	elapsedTime time.Duration
}

func sendRequest(req *http.Request) ([]byte, testData) {

	client := &http.Client{}
	time_start := time.Now()
	resp, err := client.Do(req)

	if err != nil {
		panic("aaa")
	}
	defer resp.Body.Close()
	tD := testData{
		status:      resp.Status,
		elapsedTime: time.Since(time_start),
	}

	body, err := ioutil.ReadAll(resp.Body)
	println("Respons: " + string(body))
	if err != nil {

		panic("Failed to read response")
	}

	return body, tD
}
func preformTests(testType string, amount int, start int) {
	var totalTime time.Duration
	//address := "192.168.0.119:8443" //java version
	//address := "192.168.0.119:4245" //golang version
	address := "localhost:4245" //golang version
	//change convert depending on address
	switch typeSwitch := testType; typeSwitch {
	case "query":
		for i := 0; i < amount; i++ {

			//totalTime = totalTime + tD.elapsedTime
			_, tD := sendRequest(createQueryRequest(start, address))
			totalTime = totalTime + tD.elapsedTime
			fmt.Printf("\n#%d  Status: %s   Time %s", i, tD.status, tD.elapsedTime)
			//start++
		}
	case "register":
		for i := 0; i < amount; i++ {
			_, tD := sendRequest(createRegisterRequest(start, address))
			totalTime = totalTime + tD.elapsedTime
			fmt.Printf("\n#%d  Status: %s   Time %s", i, tD.status, tD.elapsedTime)
			start++
		}

	}
	fmt.Printf("\nTotaltime %s ", totalTime)
}
func createRegisterRequest(i int, address string) *http.Request {
	serviceRegistryEntry := &ServiceRegistryEntryInput{
		ServiceDefinition: "TestSD" + strconv.Itoa(i),
		ProviderSystem: ProviderSystem{
			SystemName:         "string",
			Address:            "string",
			Port:               0,
			AuthenticationInfo: "string",
		},
		ServiceUri:    "TestSUri" + strconv.Itoa(i),
		EndOfvalidity: "2025-04-14T11:07:36.639Z",
		//EndOfvalidity: "2002-04-14T11:07:36.639Z", // test for an invalid date
		Secure: "NOT_SECURE",
		Metadata: []string{
			"metadata1",
			"metadata2",
			"metadata3",
		},

		Version: 0,
		Interfaces: []string{
			"HTTP-SECURE-JSON",
		},
	}
	body, _ := json.Marshal(serviceRegistryEntry)
	println("Sending: " + string(body))
	req, err := http.NewRequest("POST", "http://"+address+"/serviceregistry/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err.Error())
	}

	return req

}
func createQueryRequest(i int, address string) *http.Request {
	serviceQueryForm := ServiceQueryForm{
		ServiceDefinitionRequirement: "testsd" + strconv.Itoa(i),
		InterfaceRequirements:        []string{},
		SecurityRequirements:         []string{},
		MetadataRequirements:         MetadataStruct{},
		VersionRequirements:          0,
		MaxVersionRequirements:       0,
		MinVersionRequirements:       0,
		PingProviders:                false,
	}
	body, _ := json.Marshal(serviceQueryForm)
	println("Sending: " + string(body))
	req, err := http.NewRequest("POST", "http://"+address+"/serviceregistry/query", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err.Error())
	}

	return req
}
func createUnregisterRequest(i int, address string) *http.Request {
	return nil
}
