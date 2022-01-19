package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type serviceData struct {
	ID         int    `json:"id"`
	SystemName string `json:"systemName"`
	Adress     string `json:"adress"`
	Port       int    `json:"port"`
}

func query(w http.ResponseWriter, req *http.Request) {
	println(req.RemoteAddr)
	println(req.Header.Get("User-Agent"))
	switch req.Method {
	case "POST":

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		s := serviceData{}
		err := json.NewDecoder(req.Body).Decode(&s)
		s.SystemName = "newwww"
		jsonResp, err := json.Marshal(s)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(jsonResp)
		return
	default:
		println("wrong")
		fmt.Fprintf(w, "wrong method\n")

	}
}

func main() {

	http.HandleFunc("/query", query)
	http.ListenAndServe(":8090", nil)
}
