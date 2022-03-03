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
	goLocalhost = true
	printResponsAndRequest = true
	//preformTests("echo", 1, 1)
	//preformTests("register", 1, 1)
	//Testmap2()
	//preformTests("query", 1, 2) //preforms  querys
	//preformTests("unregister", 1, 1)

	//HeavyTest("register", 1000)

	amount := 1
	start := 1

	java = false //Set to true for java version
	println("Golang")
	testGoA := preformTestAll(amount, start)

	java = true
	println("Java")
	testJava := preformTestAll(amount, start)

	testGoA.printResults()

	testJava.printResults()
	//fmt.Printf("Total Difference: %s\n", testGoA.totalTime-testJava.totalTime)

}

var java bool //Set to true to convert from metadata array to struct
var goLocalhost bool
var printResponsAndRequest bool

// Preform amount test on 5 go routines
func HeavyTest(testType string, amount int) {
	go preformTests(testType, amount, 1)
	go preformTests(testType, amount, amount+1)
	go preformTests(testType, amount, (2*amount)+1)
	go preformTests(testType, amount, (3*amount)+1)
	go preformTests(testType, amount, (4*amount)+1)

	for i := 0; i < 1; {

	}
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
	if err != nil {

		panic("Failed to read response")
	}
	if printResponsAndRequest {

		if req.Body != nil {
			buf := new(bytes.Buffer)
			reqBody, _ := req.GetBody()
			buf.ReadFrom(reqBody)
			println("Sending: " + buf.String())
		}

		println("Respons: " + string(body))
	}
	// Output: Hello
	return body, tD
}
func preformTests(testType string, amount int, start int) (time.Duration, int) {
	var totalTime time.Duration
	var address string
	if java {
		address = "31.208.108.251:42454" //java version
	} else {
		if goLocalhost {
			address = "localhost:4245" //golang localversion
		} else {
			address = "31.208.108.251:4245" //golang pi version
		}

	}

	successCount := 0
	switch typeSwitch := testType; typeSwitch {
	case "query":
		println("Query")
		for i := 0; i < amount; i++ {

			//totalTime = totalTime + tD.elapsedTime
			_, tD := SendRequest(createQueryRequest(start+amount-1, address))
			totalTime = totalTime + tD.elapsedTime
			fmt.Printf("#%d   			Status: %s   			Time %s\n", i, tD.status, tD.elapsedTime)
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
			fmt.Printf("#%d   			Status: %s   			Time %s\n", i, tD.status, tD.elapsedTime)
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
			fmt.Printf("#%d   			Status: %s   			Time %s\n", i, tD.status, tD.elapsedTime)
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
			fmt.Printf("#%d   			Status: %s   			Time %s\n", i, tD.status, tD.elapsedTime)
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

		MetadataJava: map[string]string{"AAA": "aaa", "BBB": "bbb", "CCC": "ccc", "DDD": "ddd"},
		/*
			MetadataGo: []string{
				"metadata1",
				"metadata2",
				"metadata3",
			},
		*/
		Version: 0,
		Interfaces: []string{
			"HTTP-SECURE-JSON",
		},
	}
	body, _ := json.Marshal(serviceRegistryEntry)
	req, err := http.NewRequest("POST", "http://"+address+"/serviceregistry/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err.Error())
	}

	return req

}
func createQueryRequest(i int, address string) *http.Request {
	serviceQueryForm := ServiceQueryForm{
		ServiceDefinitionRequirement: "TestSD" + strconv.Itoa(i),
		InterfaceRequirements:        []string{},
		SecurityRequirements:         []string{},

		//MetadataRequirementsGo: []string{"metadata1", "metadata2", "metadata3"},
		//MetadataRequirementsGo: []string{"AAA", "BBB", "CCC", "DDD"},
		MetadataRequirementsJava: map[string]string{"AAA": "aaa", "BBB": "bbb", "CCC": "ccc", "DDD": "ddd"},
		//MetadataRequirementsJava: map[string]string{"metadata1": "", "metadata2": "", "metadata3": ""},
		VersionRequirements:    0,
		MaxVersionRequirements: 0,
		MinVersionRequirements: 0,
		PingProviders:          false,
	}
	body, _ := json.Marshal(serviceQueryForm)

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
	fmt.Printf("Echo total time: %s, 				Succes count: %d \n", cTR.echoTime, cTR.echoSuccesCount)
	fmt.Printf("Query total time: %s, 				Succes count: %d \n", cTR.queryTime, cTR.querySuccesCount)
	fmt.Printf("Register total time: %s, 			Success count: %d \n", cTR.registerTime, cTR.registerSuccesCount)
	fmt.Printf("Unregister total time: %s, 			success count: %d \n", cTR.unregisterTime, cTR.unregisterSuccesCount)
	fmt.Printf("Combined total time: %s, 			Total success count: %d \n", cTR.totalTime, cTR.totalSuccesCount)
}
