package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

// OpenDatabase is functions to open database connectivity
func OpenDatabase(dsn string) {
	var err error
	db, err = sql.Open("sqlite3", dsn)

	handleError("could not open database: %v", err)

	if err := db.Ping(); err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
}

func handleError(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

// Register model.
type Register struct {
	ID                 uint   `json:"id"`
	ServiceDefinition  string `json: "serviceDefinition"`
	SystemName         string `json:"systemName"`
	Port               string `json: "port"`
	AuthenticationInfo string `json:"authenticationInfo"`
	ServiceURI         string `json:"serviceURI"`
	EndOfvalidity      string `json:"endOfValidity"`
	Secure             string `json: NOT_SECURE`
	Address            string `json:"address"`
	MetaData           string `json:"metaData"`
	Version            string `json: "version"`
	Interfaces         string `json:"interfaces"`
}

// All functions to returns all Service list.
func (model *Register) All() []Register {
	rows, err := db.Query("SELECT * FROM Services")
	handleError("query failed: %v", err)

	var Services []Register
	defer rows.Close()

	for rows.Next() {
		var Service Register

		err := rows.Scan(&Service.ID, &Service.SystemName, &Service.MetaData)
		handleError("scan failed: %v", err)

		Services = append(Services, Service)
	}

	return Services
}

// DATABASE HANDLER
func (model *Register) Save() *Register {
	stmt, err := db.Prepare("INSERT INTO Services (serviceDefinition, SystemName, metaData, Port, authenticationInfo, serviceURI, endOfValidity, secure, address, version, interfaces) VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	handleError("could prepare statement: %v", err)
	println(model.Port, "hello")
	println(model.SystemName, "hello2")
	res, err := stmt.Exec(model.ServiceDefinition, model.SystemName, model.MetaData, model.Port, model.AuthenticationInfo, model.ServiceURI, model.EndOfvalidity, model.Secure, model.Address, model.Version, model.Interfaces)
	handleError("failed to store: %v", err)

	id, _ := res.LastInsertId()
	model.ID = uint(id)
	println("store worked")
	return model
}

// Find Service by id.
func (model *Register) Find(name string) *Register {
	row := db.QueryRow("SELECT * FROM Services WHERE SystemName = ?", name)
	println("?", name)
	Service := &Register{}
	err := row.Scan(&Service.SystemName)

	if err != nil {
		if err == sql.ErrNoRows {
			println("ERROR with query")
			return nil
		}
	}

	return Service
}

func (model *Register) ToggleDoneStatus() bool {

	stmt, err := db.Prepare("UPDATE Services SET metaData = ? WHERE id = ?")
	handleError("could not prepare statement: %v", err)

	res, err := stmt.Exec(model.MetaData, model.ID)
	handleError("query failed: %v", err)

	affecteds, err := res.RowsAffected()
	handleError("could not get affected rows: %v", err)

	if affecteds > 0 {
		return true
	}

	return false
}

// Delete is functions to remove Service from database.
func (model *Register) Delete() bool {
	stmt, err := db.Prepare("DELETE  FROM Services WHERE SystemName = ?")
	handleError("could not prepare statement: %v", err)

	res, err := stmt.Exec(model.SystemName)
	handleError("query failed: %v", err)

	affecteds, err := res.RowsAffected()
	handleError("could not get affected rows: %v", err)
	check := affecteds
	if check > 0 {
		return true
	}

	return false
}
