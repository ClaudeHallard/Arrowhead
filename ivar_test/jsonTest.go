package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//type serviceData struct {
//ID         int    `json:"id"`
//SystemName string `json:"systemName"`
//Adress     string `json:"adress"`
//Port       int    `json:"port"`
//}

func createJSON() {
	sD := serviceData{
		ID:         7,
		SystemName: "7th",
		Adress:     "7thAdr",
		Port:       77,
	}
	dat, err := json.Marshal(sD)
	if err != nil {

	}
	err = ioutil.WriteFile("./database/file.json", dat, 0644)
	if err != nil {

	}
}
func main() {
	dat, err := ioutil.ReadFile("./database/file.json")
	if err != nil {

	}
	println(string(dat))
	sD := serviceData{}
	err = json.Unmarshal(dat, &sD)
	fmt.Println(sD)
	if err != nil {

	}
}
