package main

import (
	"database/sql"
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

func (model *ServiceQueryForm) Query() *ServiceQueryList {
	serviceQueryList := &ServiceQueryList{}
	serviceQueryList.ServiceQueryData = getServiceByID(-1)
	serviceQueryList.UnfilteredHits = len(serviceQueryList.ServiceQueryData)
	serviceQueryList.serviceDefenitionFilter(*model)

	return serviceQueryList
}

// check for serviceURI AND / OR serviceDefinition
func GetCount(servicedef string, serviceURI string) int {
	var cnt int
	//_ = db.QueryRow(`select count(*) from Services where serviceDefinition= ? AND serviceURI = ?`, servicedef, serviceURI).Scan(&cnt)
	_ = db.QueryRow(`select count(*) from Services where serviceURI = ?`, serviceURI).Scan(&cnt)
	return cnt
}

// function to store the register input
func (model *ServiceRegistryEntryInput) Save() *ServiceRegistryEntryOutput {

	// check if the service exist
	if GetCount(model.ServiceDefinition, model.ServiceUri) > 0 {

	} else {

		stmt, err := db.Prepare("INSERT INTO Services (serviceDefinition, systemName, address, port, authenticationInfo, serviceURI, endOfValidity, secure, version, createdAt, updatedAt) VALUES (?,?,?,?,?,?,?,?,?,?,?)")

		currentTime := time.Now().Format("2006-01-02 15:04:05")

		res, err := stmt.Exec(model.ServiceDefinition, model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port, model.ProviderSystem.AuthenticationInfo, model.ServiceUri, model.EndOfvalidity, model.Secure, model.Version, currentTime, currentTime)
		if err != nil {
			println(err.Error())
			panic("Encounterd an error during registration while inserting service")
		}
		lastId, err := res.LastInsertId()
		handleError("failed to store: %v", err)
		for _, v := range model.MetadataGo {
			stmt, err := db.Prepare("INSERT INTO MetaData (serviceID, metaData) VALUES (?,?)")
			if err != nil {
				println(err.Error())
				panic("Encounterd an error during registration while inserting metadata")
			}
			_, err = stmt.Exec(lastId, v)

		}

		//loop through the interfaces array and add them to the table.
		for _, v := range model.Interfaces {
			stmt, err := db.Prepare("INSERT INTO Interfaces (serviceID, interfaceName, createdAt, updatedAt) VALUES (?,?,?,?)")
			handleError("could prepare statement: %v", err)
			_, err = stmt.Exec(lastId, v, currentTime, currentTime)
		}
		handleError("failed to store: %v", err)

		return &getServiceByID(lastId)[0]
	}
	return nil
}

// function to store the register input. Looping through struct instead of string array.
func (model *ServiceRegistryEntryInput) SaveJava() *ServiceRegistryEntryOutput {

	// check if the service exist
	if GetCount(model.ServiceDefinition, model.ServiceUri) > 0 {

	} else {

		stmt, err := db.Prepare("INSERT INTO Services (serviceDefinition, systemName, address, port, authenticationInfo, serviceURI, endOfValidity, secure, version, createdAt, updatedAt) VALUES (?,?,?,?,?,?,?,?,?,?,?)")
		currentTime := time.Now().Format("2006-01-02 15:04:05")

		res, err := stmt.Exec(model.ServiceDefinition, model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port, model.ProviderSystem.AuthenticationInfo, model.ServiceUri, model.EndOfvalidity, model.Secure, model.Version, currentTime, currentTime)
		println("hit")
		if err != nil {
			println(err.Error())
			panic("Encounterd an error during registration while inserting service")
		}
		lastId, err := res.LastInsertId()
		handleError("failed to store: %v", err)

		// Insert metadata
		stmt1, err := db.Prepare("INSERT INTO MetaData (serviceID, metaData) VALUES (?,?)")
		metadata1 := model.MetadataJava.AdditionalProp1
		println(model.MetadataJava.AdditionalProp1)
		print("hej")
		if len(metadata1) == 0 {
			println("No metadata1")
		} else {
			_, err = stmt1.Exec(lastId, metadata1)

			stmt2, err := db.Prepare("INSERT INTO MetaData (serviceID, metaData) VALUES (?,?)")
			metadata2 := model.MetadataJava.AdditionalProp2
			if len(metadata2) > 0 && err == nil {
				_, err = stmt2.Exec(lastId, metadata2)
			} else {
				println("No metadata2")
			}

			stmt3, err := db.Prepare("INSERT INTO MetaData (serviceID, metaData) VALUES (?,?)")
			metadata3 := model.MetadataJava.AdditionalProp3
			if len(metadata3) > 0 && err == nil {
				_, err = stmt3.Exec(lastId, metadata2)
			} else {
				println("No metadata3")
			}
		}

		//loop through the interfaces array and add them to the table.
		for _, v := range model.Interfaces {
			stmt, err := db.Prepare("INSERT INTO Interfaces (serviceID, interfaceName, createdAt, updatedAt) VALUES (?,?,?,?)")
			handleError("could prepare statement: %v", err)
			_, err = stmt.Exec(lastId, v, currentTime, currentTime)
		}
		handleError("failed to store: %v", err)

		return &getServiceByID(lastId)[0]
	}
	return nil
}

// Function to delete a service and the assosiated data
func (model *ServiceRegistryEntryInput) Delete() bool {
	stmt, err := db.Prepare("DELETE FROM Services WHERE serviceDefinition = ? AND systemName = ? AND  address = ? AND port = ?")
	handleError("could not prepare statement: %v", err)
	res, err := stmt.Exec(model.ServiceDefinition, model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port)
	handleError("query failed: %v", err)
	affecteds, err := res.RowsAffected()
	handleError("could not get affected rows: %v", err)
	check := affecteds
	if check > 0 {
		return true
	}
	return false
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

func getServiceByID(id int64) []ServiceRegistryEntryOutput { //If id <= 0 then all services will be retured
	var serviceList []ServiceRegistryEntryOutput
	//Get the interfaces
	var interfaces = getInterfaceByID(id)

	//Get the metadata
	var metaServiceID, metaData = getMetadataByID(id)

	var rows *sql.Rows
	var err error
	if 0 <= id {
		rows, err = db.Query("SELECT * FROM Services WHERE id=?", id)
	} else {
		rows, err = db.Query("SELECT * FROM Services")
	}

	//Get the services and add metadata and interfaces to their corresponding service
	handleError("query failed: %v", err)
	for rows.Next() {

		var service ServiceRegistryEntryOutput

		err := rows.Scan(&service.ID, &service.ServiceDefinition.ServiceDefinition, &service.ServiceDefinition.CreatedAt, &service.ServiceDefinition.UpdatedAt,
			&service.Provider.SystemName, &service.Provider.Address, &service.Provider.Port, &service.Provider.AuthenticationInfo, &service.Provider.CreatedAt, &service.Provider.UpdatedAt,
			&service.ServiceUri, &service.EndOfValidity, &service.Secure, &service.Version, &service.CreatedAt, &service.UpdatedAt)
		//add things to other fields

		if err != nil {
			println(err.Error())
			panic("Encounterd an error while quering the database for services")
		}
		for i := 0; i < len(metaData); i++ {

			if service.ID == metaServiceID[i] {
				service.Metadata = append(service.Metadata, metaData[i])
			}
		}
		for _, v := range interfaces {
			if service.ID == v.ID {
				service.Interfaces = append(service.Interfaces, v)
			}
		}

		serviceList = append(serviceList, service)

	}

	defer rows.Close()
	return serviceList
}

func getInterfaceByID(id int64) []Interface { //If id <= 0  then all interfaces will be retured
	var interfaces []Interface
	var rows *sql.Rows
	var err error
	if 0 <= id {
		rows, err = db.Query("SELECT serviceID, interfaceName, createdAt, updatedAt FROM Interfaces WHERE serviceID = ?", id)
	} else {
		rows, err = db.Query("SELECT serviceID, interfaceName, createdAt, updatedAt FROM Interfaces")
	}

	//rows, err := db.Query("SELECT serviceID, interfaceName, createdAt, updatedAt FROM Interfaces WHERE serviceID = ?", id)
	handleError("query failed: %v", err)
	for rows.Next() {
		var i Interface
		err := rows.Scan(&i.ID, &i.InterfaceName, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			println(err.Error())
			panic("Encounterd an error while scaning database for interfaces")
		}
		interfaces = append(interfaces, i)

	}
	defer rows.Close()
	return interfaces
}

func getMetadataByID(id int64) ([]int, []string) { //If id <= 0 then all metadata will be retured
	var serviceID []int
	var metaData []string
	var rows *sql.Rows
	var err error
	if 0 <= id {
		rows, err = db.Query("SELECT serviceID, metaData  FROM MetaData WHERE serviceID = ?", id)
	} else {
		rows, err = db.Query("SELECT serviceID, metaData  FROM MetaData")
	}

	handleError("query failed: %v", err)
	for rows.Next() {
		var sID int
		var mData string

		err := rows.Scan(&sID, &mData)

		if err != nil {
			println(err.Error())
			panic("Encounterd an error while scaning database for metadata")
		}

		serviceID = append(serviceID, sID)
		metaData = append(metaData, mData)

	}
	defer rows.Close()
	return serviceID, metaData
}
func deleteByID(id int) {
	stmt, err := db.Prepare("DELETE FROM Services WHERE id = ?")
	handleError("could not prepare statement: %v", err)
	_, err = stmt.Exec(id)

	handleError("delete failed: %v", err)

}
func cleanPastValidityDate() {
	var pastValidityServices []int
	var id int
	var endOfValidity string
	rows, err := db.Query("SELECT id, endoOfValidity FROM Services")
	handleError("query failed: %v", err)
	for rows.Next() {
		err := rows.Scan(&id, &endOfValidity)
		handleError("query failed: %v", err)
		if !validityCheck(endOfValidity) {
			pastValidityServices = append(pastValidityServices, id)
		}
	}
	defer rows.Close()
	for _, v := range pastValidityServices {
		deleteByID(v)
	}

}
func validityCheck(timeString string) bool {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {

		println(err.Error())

		return false
	}
	return time.Now().Before(t)

}
