package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type serviceData struct {
	ID         int    `json:"id"`
	SystemName string `json:"systemName"`
	Adress     string `json:"adress"`
	Port       int    `json:"port"`
}

func main() {
	sD := &serviceData{
		ID:         7,
		SystemName: "7th",
		Adress:     "7thAdr",
		Port:       77,
	}
	println("Sending: ", sD.ID, sD.SystemName, sD.Adress, sD.Port)
	body, _ := json.Marshal(sD)
	resp, err := http.Post("http://127.0.0.1:8090/query", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//Check response code, (201)
	if resp.StatusCode == http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
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
