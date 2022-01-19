package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main2() {

	resp, err := http.Get("http://127.0.0.1:4245/serviceregistry/echo")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
