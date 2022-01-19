package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Unregister struct {
	systemName string `json:"systemName"`
}

func main() {
	sD := &Unregister{
		systemName: "bb",
	}
	json_data, err := json.Marshal(sD)
	/**
	values := map[string]string{
		systemName: "bb"}
	fmt.Println("Performing Http Delete...")
	json_data, err := json.Marshal(values)
	**/

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodDelete, "http://127.0.0.1:4245/serviceregistry/unregister", bytes.NewBuffer(json_data))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to string
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	/**
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res["json"])
	**/
}
