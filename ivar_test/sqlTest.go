package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	insertJSON()
}
func searchSQL() {

	db, err := sql.Open("sqlite3", "./database/go.db")
	if err != nil {
		panic(err)
	}
	searchID := 2
	var name, adress string
	var port int
	sqlStatement := "SELECT systemName, adress, port FROM Service WHERE id=?"

	row := db.QueryRow(sqlStatement, searchID)
	err = row.Scan(&name, &adress, &port)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			panic(err)
		}
		fmt.Println(name)
	} else {
		fmt.Println(name, adress, port)
	}

	//dont know, dont touch
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}
func insertJSON() {
	db, err := sql.Open("sqlite3", "./database/go.db")
	dat, err := ioutil.ReadFile("./database/file.json")
	sqlStatement := "INSERT INTO Service (systemName, adress, port) VALUES ($1, $2, $3)"
	if err != nil {

	}
	sD := serviceData{}
	err = json.Unmarshal(dat, &sD)
	_, err = db.Exec(sqlStatement, sD.SystemName, sD.Adress, sD.Port)
	if err != nil {
		panic(err)
	}

}
