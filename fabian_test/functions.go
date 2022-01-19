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

type Query struct {
	ServiceDefinitionRequirement string `json: serviceDefinitionRequirement`
	InterfaceRequirements        string `json: interfaceRequirements`
	SecurityRequirements         string `json: securityRequirements`
	MetadataRequirements         string `json: metadataRequirements`
	VersionRequirement           uint   `json: versionRequirement`
	MaxVersionRequirement        uint   `json: maxVersionRequirement`
	MinVersionRequirement        uint   `json: minVersionRequirement`
	PingProviders                bool   `json: pingProviders`
}

// function to return the specific service
func (model *Register) Query(name string) []Register {
	rows, err := db.Query("SELECT * FROM Services WHERE systemName = ?", name)
	handleError("query failed: %v", err)

	var Services []Register
	defer rows.Close()

	for rows.Next() {
		var Service Register

		err := rows.Scan(&Service.ID, &Service.ServiceDefinition, &Service.SystemName, &Service.Port, &Service.AuthenticationInfo, &Service.ServiceURI, &Service.EndOfvalidity, &Service.Secure, &Service.Address, &Service.MetaData, &Service.Version, &Service.Interfaces)
		handleError("scan failed: %v", err)

		Services = append(Services, Service)
	}

	return Services
}

// function to store the register input
func (model *ServiceRegistryEntryInput) Save() *ServiceRegistryEntryInput {

	println(model.ServiceDefinition)
	stmt, err := db.Prepare("INSERT INTO Services (serviceDefinition, systemName, address, port, authenticationInfo, serviceUri, endOfValidity, secure, metaData1, metaData2, metaData3, version, createdAt, updatedAt) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	handleError("could prepare statement: %v", err)
	fmt.Println("Registration recived", model.ProviderSystem.SystemName)
	res, err := stmt.Exec(model.ServiceDefinition, model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port, model.ProviderSystem.AuthenticationInfo, model.ServiceUri, model.EndOfvalidity, model.Secure, model.Metadata.AdditionalProp1, model.Metadata.AdditionalProp2, model.Metadata.AdditionalProp3, model.Version, currentTime, currentTime)
	lastId, err := res.LastInsertId()
	handleError("failed to store: %v", err)

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
