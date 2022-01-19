package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
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
	ServiceDefinition  string `json:"serviceDefinition"`
	SystemName         string `json:"systemName"`
	Port               string `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
	ServiceURI         string `json:"serviceURI"`
	EndOfvalidity      string `json:"endOfValidity"`
	Secure             string `json: NOT_SECURE`
	Address            string `json:"address"`
	MetaData           string `json:"metaData"`
	Version            string `json: "version"`
	Interfaces         string `json:"interfaces"`
}

// function to return the specific service
func (model *ServiceQueryForm) Query() *ServiceQueryList {
	//serviceQueryList := &ServiceQueryList{}
	serviceQueryList := getAllServices()
	serviceQueryList.serviceDefenitionFilter(*model)

	return serviceQueryList
}

// function to store the register input
func (model *ServiceRegistryEntryInput) Save() *ServiceRegistryEntryInput {

	println(model.ServiceDefinition)
	stmt, err := db.Prepare("INSERT INTO Services (serviceDefinition, systemName, address, port, authenticationInfo, serviceURI, endOfValidity, secure, version, createdAt, updatedAt) VALUES (?,?,?,?,?,?,?,?,?,?,?)")

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	handleError("could prepare statement: %v", err)
	fmt.Println("Registration recived", model.ProviderSystem.SystemName)
	res, err := stmt.Exec(model.ServiceDefinition, model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port, model.ProviderSystem.AuthenticationInfo, model.ServiceUri, model.EndOfvalidity, model.Secure, model.Version, currentTime, currentTime)
	lastId, err := res.LastInsertId()
	handleError("failed to store: %v", err)
	for _, v := range model.Metadata {
		stmt, err := db.Prepare("INSERT INTO MetaData (serviceID, metaData) VALUES (?,?)")
		handleError("could prepare statement: %v", err)
		_, err = stmt.Exec(lastId, v)

	}
	//loop through the interfaces array and add them to the table.
	for _, v := range model.Interfaces {
		stmt, err := db.Prepare("INSERT INTO Interfaces (serviceID, interfaceName, createdAt, updatedAt) VALUES (?,?,?,?)")
		handleError("could prepare statement: %v", err)
		_, err = stmt.Exec(lastId, v, currentTime, currentTime)
	}
	handleError("failed to store: %v", err)

	println("LastInsertId: ", lastId)
	println(len(model.Interfaces))
	//id, _ := res.LastInsertId()
	//model.ID = uint(id)

	return model
}

// function to find Service.
func (model *Register) Find(name string) *Register {
	row := db.QueryRow("SELECT * FROM Services WHERE SystemName = ?", name)

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

// function to delete service from database.
func (model *Register) Delete(name string) bool {
	stmt, err := db.Prepare("DELETE FROM Services WHERE SystemName = ?")
	handleError("could not prepare statement: %v", err)
	res, err := stmt.Exec(name)
	handleError("query failed: %v", err)

	affecteds, err := res.RowsAffected()
	handleError("could not get affected rows: %v", err)
	check := affecteds
	if check > 0 {
		return true
	}
	return false
}
func getAllServices() *ServiceQueryList {
	var interfaces []Interface
	var serviceQueryList = &ServiceQueryList{}
	rows, err := db.Query("SELECT serviceID, interfaceName, createdAt, updatedAt FROM Interfaces")
	handleError("query failed: %v", err)
	for rows.Next() {
		var i Interface
		err := rows.Scan(&i.ID, &i.InterfaceName, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			handleError("scan failed: %v", err)
		}
		interfaces = append(interfaces, i)

	}
	defer rows.Close()
	rows, err = db.Query("SELECT * FROM Services")
	handleError("query failed: %v", err)

	for rows.Next() {

		var service ServiceRegistryEntryOutput

		err := rows.Scan(&service.ID, &service.ServiceDefinition.ServiceDefinition, &service.ServiceDefinition.CreatedAt, &service.ServiceDefinition.UpdatedAt,
			&service.Provider.SystemName, &service.Provider.Address, &service.Provider.Port, &service.Provider.AuthenticationInfo, &service.Provider.CreatedAt, &service.Provider.UpdatedAt,
			&service.ServiceUri, &service.EndOfValidity, &service.Secure, &service.Version, &service.CreatedAt, &service.UpdatedAt)
		//add things to other fields

		if err != nil {
			println(err.Error())

			panic("Check for null values in the database")
			handleError("scan failed: %v", err)
		}

		for _, v := range interfaces {
			if service.ID == v.ID {
				service.Interfaces = append(service.Interfaces, v)
			}
		}
		serviceQueryList.ServiceQueryData = append(serviceQueryList.ServiceQueryData, service)
		serviceQueryList.UnfilteredHits = serviceQueryList.UnfilteredHits + 1
	}
	defer rows.Close()

	return serviceQueryList

}

func (serviceQueryList *ServiceQueryList) serviceDefenitionFilter(serviceQueryForm ServiceQueryForm) *ServiceQueryList {
	var queryHits []ServiceRegistryEntryOutput
	for _, v := range serviceQueryList.ServiceQueryData {
		if v.ServiceDefinition.ServiceDefinition == serviceQueryForm.ServiceDefinitionRequirement {
			queryHits = append(queryHits, v)
		}

	}
	serviceQueryList.ServiceQueryData = queryHits
	return serviceQueryList

}
