package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	java = false //Set to true for java version
	//preformTests("echo", 100, 0)
	//preformTests("register", 30, 1)
	//preformTests("unregister", 3, 1)
	//preformTests("echo", 100, 0) // preform 100 echo calls
	//preformTests("query", 10000, 40) //preforms 100 querys
	//_, td := SendRequest(createUnregisterRequestJava(101, "31.208.108.251:42454")
	println("Golang")
	amount := 1000
	start := 1
	testResult1 := preformTestAll(amount, start)
	println("Java")
	java = true
	testResult2 := preformTestAll(amount, start)
	testResult1.printResults()
	testResult2.printResults()

	//createUnregisterRequestJava()
	//preformTests("register", 5, 110) //registers 5 services, starting at index 100
	//preformTests("unregister", 20, 90) // unregisters 5 services, starting at index 100

}

var java bool //Set to true to convert from metadata array to struct

// asdasdasd
func (metadata Metadata) MarshalJSON() ([]byte, error) {
	metadataStruct := MetadataStruct{}
	if java {
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

// To do that, simply add an extra indent to your comment's text.
func SendRequest(req *http.Request) ([]byte, testData) {

	client := &http.Client{}
	time_start := time.Now()
	resp, err := client.Do(req)

	if err != nil {
		panic(err.Error())
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
	// Output: Hello
	return body, tD
}
func preformTests(testType string, amount int, start int) (time.Duration, int) {
	var totalTime time.Duration
	var address string
	if java {
		address = "31.208.108.251:42454"
	} else {
		//address = "31.208.108.251:4245"
		address = "localhost:4245"
	}
	//address := "31.208.108.251:42454" //java version
	//address = "31.208.108.251:4245" //golang version
	//address := "localhost:4245" //golang version
	successCount := 0
	switch typeSwitch := testType; typeSwitch {
	case "query":
		println("Query")
		for i := 0; i < amount; i++ {

			//totalTime = totalTime + tD.elapsedTime
			_, tD := SendRequest(createQueryRequest(start, address))
			totalTime = totalTime + tD.elapsedTime
			fmt.Printf("#%d  Status: %s   Time %s\n", i, tD.status, tD.elapsedTime)
			if strings.Contains(tD.status, "200") {
				successCount = successCount + 1
			}
			//start++
		}

	case "register":
		println("Register")
		for i := 0; i < amount; i++ {
			_, tD := SendRequest(createRegisterRequest(start, address))
			totalTime = totalTime + tD.elapsedTime
			fmt.Printf("#%d  Status: %s   Time %s\n", i, tD.status, tD.elapsedTime)
			start++
			if strings.Contains(tD.status, "201") {
				successCount = successCount + 1
			}
		}

	case "unregister":
		println("Unregister")
		for i := 0; i < amount; i++ {
			_, tD := SendRequest(createUnregisterRequest(start, address))

			totalTime = totalTime + tD.elapsedTime
			fmt.Printf("#%d  Status: %s   Time %s\n", i, tD.status, tD.elapsedTime)
			start++
			if strings.Contains(tD.status, "200") {
				successCount = successCount + 1
			}
		}

	case "echo":
		println("Echo")
		for i := 0; i < amount; i++ {
			_, tD := SendRequest(echoRequest(start, address))
			totalTime = totalTime + tD.elapsedTime
			fmt.Printf("#%d  Status: %s   Time %s\n", i, tD.status, tD.elapsedTime)
			start++
			if strings.Contains(tD.status, "200") {

				successCount = successCount + 1
			}
		}
	}

	fmt.Printf("Totaltime %s \n", totalTime)
	return totalTime, successCount
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
		//ServiceDefinitionRequirement: "TestSD" + strconv.Itoa(i),
		ServiceDefinitionRequirement: "TestSD101",
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
	req, _ := http.NewRequest("DELETE", "http://"+address+"/serviceregistry/unregister", nil)
	q := req.URL.Query()
	q.Add("address", "string")
	q.Add("port", "0")
	q.Add("service_definition", "TestSD"+strconv.Itoa(i))
	q.Add("service_uri", "TestSUri"+strconv.Itoa(i))
	q.Add("system_name", "string")
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	//client := &http.Client{}
	//println(r.Header)
	//resp, _ := client.Do(req)
	//fmt.Println(resp.Status)
	return req
}

func echoRequest(i int, address string) *http.Request {

	req, err := http.NewRequest("GET", "http://"+address+"/serviceregistry/echo", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err.Error())
	}
	return req
}
func preformTestAll(amount int, start int) combinedTestResult {
	result := combinedTestResult{}
	result.java = java
	result.echoTime, result.echoSuccesCount = preformTests("echo", amount, start)
	result.registerTime, result.registerSuccesCount = preformTests("register", amount, start)
	result.queryTime, result.querySuccesCount = preformTests("query", amount, start)
	result.unregisterTime, result.unregisterSuccesCount = preformTests("unregister", amount, start)
	result.totalTime = result.echoTime + result.queryTime + result.registerTime + result.unregisterTime
	result.totalSuccesCount = result.echoSuccesCount + result.querySuccesCount + result.registerSuccesCount + result.unregisterSuccesCount
	return result
}

type combinedTestResult struct {
	java     bool
	echoTime time.Duration

	queryTime             time.Duration
	registerTime          time.Duration
	unregisterTime        time.Duration
	totalTime             time.Duration
	echoSuccesCount       int
	querySuccesCount      int
	registerSuccesCount   int
	unregisterSuccesCount int
	totalSuccesCount      int
}

func (cTR combinedTestResult) printResults() {
	fmt.Printf("Java version: %t\n", cTR.java)
	fmt.Printf("Echo total time: %s, 			succes count: %d \n", cTR.echoTime, cTR.echoSuccesCount)
	fmt.Printf("Query total time: %s, 			succes count: %d \n", cTR.queryTime, cTR.querySuccesCount)
	fmt.Printf("Register total time: %s, 		succes count: %d \n", cTR.registerTime, cTR.registerSuccesCount)
	fmt.Printf("Unregister total time: %s, 		succes count: %d \n", cTR.unregisterTime, cTR.unregisterSuccesCount)
	fmt.Printf("Combined total time: %s, 		total succes count: %d \n", cTR.totalTime, cTR.totalSuccesCount)
}
