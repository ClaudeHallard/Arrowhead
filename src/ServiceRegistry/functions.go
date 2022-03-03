/*
**********************************************
@authors:
Ivar Grunaeu, ivagru-9@student.ltu.se
Jean-Claude Hallard, jeahal-8@student.ltu.se
Marcus Paulsson, marpau-8@student.ltu.se
Pontus Schünemann, ponsch-9@student.ltu.se
Fabian Widell, fabwid-9@student.ltu.se

**********************************************
functions.go contains all functions for communicating with
the database, these functions are called and used by methods
in other files. It also contains numerous helper functions.  */

package ServiceRegistry

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

var db *sql.DB

// OpenDatabase establishing the connection to the specified database file.
func OpenDatabase() {

	var err error
	db, err = sql.Open("sqlite3", "file:ServiceRegistry/database/registryDB.db?cache=shared&_journal=WAL&_foreign_keys=on")
	if err != nil {
		panic(err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

}

// Helper function for printing error messages.
func handleError(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

// Query invokes checks on the query form and sends request to extract information from the database. Returns the information in a query list.
func (model *ServiceQueryForm) Query() *ServiceQueryList {

	if len(model.MetadataRequirementsGo) != 0 {

		model.MetadataRequirementsJava = convertToJavaMetadata(model.MetadataRequirementsGo)

	}

	serviceQueryList := &ServiceQueryList{}
	serviceQueryList.ServiceDefenitionFilter(*model)

	if int(len(model.MetadataRequirementsJava)) > 0 {

		serviceQueryList.metadataRequiermentFilter(*model)
	}
	if len(model.MetadataRequirementsGo) != 0 {
		for i := 0; i < len(serviceQueryList.ServiceQueryData); i++ {
			serviceQueryList.ServiceQueryData[i].MetadataGo = convertToGoMetadata(serviceQueryList.ServiceQueryData[i].MetadataJava)
		}
	}

	return serviceQueryList
}

// Update function that is invoked if the service allready exist when register, then the service
// is updated instead of registed again. Returns a request for the updated information
func (model *ServiceRegistryEntryInput) updateService() *ServiceRegistryEntryOutput {

	stmt, err := db.Prepare("UPDATE Services SET (systemName, address, port, authenticationInfo, endOfValidity, secure, version, updatedAt) = (?,?,?,?,?,?,?,?) WHERE serviceDefinition = ? AND serviceURI = ?")

	currentTime := time.Now().Format("2006-01-02T15:04:051Z")
	_, err = stmt.Exec(model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port, model.ProviderSystem.AuthenticationInfo, model.EndOfvalidity, model.Secure, model.Version, currentTime, model.ServiceDefinition, model.ServiceUri)
	defer stmt.Close()
	if err != nil {
		println(err.Error())
		panic("Encounterd an error during registration while inserting service, (updateService)")
	}

	var cnt int64
	row, err := db.Query(`select id from Services where serviceURI = ?`, model.ServiceUri)
	defer row.Close()
	if err != nil {
		panic("Encounterd an error during updateService" + err.Error())
	}
	row.Scan(&cnt)
	return &getServiceByID(cnt)[0]
}

// Check for serviceURI AND serviceDefinition, if a postive value is returned it means
// a similar service exist in the service registry.
func GetCountUniqueURI(serviceURI string) int {
	var cnt int
	row, err := db.Query(`select * from Services where serviceURI= ?`, serviceURI)
	defer row.Close()
	row.Scan(&cnt)
	if err != nil {
		panic("Encounterd an error during GetCountUniqueURI" + err.Error())
	}

	return cnt
}

// Check for serviceURI AND serviceDefinition if the service shall be updated.
func GetCountUpdate(servicedef string, serviceURI string) int {
	var cnt int
	row, err := db.Query(`select count(*) from Services where serviceURI = ? AND serviceDefinition = ?`, serviceURI, servicedef)
	defer row.Close()
	row.Scan(&cnt)
	if err != nil {
		panic("Encounterd an error during GetCountUpdate" + err.Error())
	}
	return cnt
}

// Function to store the register input into the database, returns the confirmation form if
// it succeeded and a null value if it didnt.
func (model *ServiceRegistryEntryInput) Save() *ServiceRegistryEntryOutput {

	if !validityCheck(model.EndOfvalidity) {
		return nil
	}
	if len(model.MetadataGo) != 0 {
		model.MetadataJava = convertToJavaMetadata(model.MetadataGo)
	}

	// check if the service exist in database
	if GetCountUniqueURI(model.ServiceUri) > 0 {
		if GetCountUpdate(model.ServiceDefinition, model.ServiceUri) > 0 {
			println("Register worked, service updated")
			return model.updateService() // updates the service in the database

		}
		return nil

	} else { // create service if it not exist.

		stmt, err := db.Prepare("INSERT INTO Services (serviceDefinition, systemName, address, port, authenticationInfo, serviceURI, endOfValidity, secure, version, createdAt, updatedAt) VALUES (?,?,?,?,?,?,?,?,?,?,?)")

		currentTime := time.Now().Format("2006-01-02T15:04:05.000Z")

		res, err := stmt.Exec(model.ServiceDefinition, model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port, model.ProviderSystem.AuthenticationInfo, model.ServiceUri, model.EndOfvalidity, model.Secure, model.Version, currentTime, currentTime)
		defer stmt.Close()
		if err != nil {
			println(err.Error())
			panic("Encounterd an error during registration while inserting service, (Save)")
		}
		lastId, err := res.LastInsertId()

		handleError("failed to store: %v", err)

		// Loop through the metadata array and add them to the table.
		for key, value := range model.MetadataJava {
			stmt, err := db.Prepare("INSERT INTO Metadata (serviceID, key, value) VALUES (?, ?, ?)")
			if err != nil {
				println(err.Error())
				panic("Encounterd an error during registration while inserting Metadata")
			}
			_, err = stmt.Exec(lastId, key, value)
			defer stmt.Close()
		}

		// Loop through the interfaces array and add them to the table.
		for _, v := range model.Interfaces {
			stmt, err := db.Prepare("INSERT INTO Interfaces (serviceID, interfaceName, createdAt, updatedAt) VALUES (?,?,?,?)")
			handleError("could prepare statement: %v", err)
			_, err = stmt.Exec(lastId, v, currentTime, currentTime)
			defer stmt.Close()
		}

		handleError("failed to store: %v", err)
		println("Register worked")

		// Setup the return form with all the params if registration succeeded.
		returnForm := ServiceRegistryEntryOutput{
			ID: int(lastId),
			ServiceDefinition: ServiceDefinition{
				ID:                int(lastId),
				ServiceDefinition: model.ServiceDefinition,
				CreatedAt:         currentTime,
				UpdatedAt:         currentTime,
			},
			Provider: Provider{
				ID:                 int(lastId),
				SystemName:         model.ProviderSystem.SystemName,
				Address:            model.ProviderSystem.Address,
				Port:               model.ProviderSystem.Port,
				AuthenticationInfo: model.ProviderSystem.AuthenticationInfo,
				CreatedAt:          currentTime,
				UpdatedAt:          currentTime,
			},
			ServiceUri:    model.ServiceUri,
			EndOfValidity: model.EndOfvalidity,
			Secure:        model.Secure,
			MetadataGo:    model.MetadataGo,
			MetadataJava:  model.MetadataJava,
			Version:       model.Version,
			CreatedAt:     currentTime,
			UpdatedAt:     currentTime,
		}

		interfaceArray := make([]Interface, len(model.Interfaces))
		for i := 0; i < len(model.Interfaces); i++ {
			interfaceArray[i] = Interface{
				ID:            int(lastId),
				InterfaceName: model.Interfaces[i],
				CreatedAt:     currentTime,
				UpdatedAt:     currentTime,
			}
		}
		returnForm.Interfaces = interfaceArray
		if len(model.MetadataGo) != 0 {
			returnForm.MetadataJava = convertToJavaMetadata(model.MetadataGo)
		}

		return &returnForm
	}

}

// Function to delete a service and the assosiated data.
func (model *ServiceRegistryEntryInput) Delete() bool {

	stmt, err := db.Prepare("DELETE FROM Services WHERE serviceDefinition = ? AND systemName = ? AND  address = ? AND port = ?")
	handleError("could not prepare statement: %v", err)
	res, err := stmt.Exec(model.ServiceDefinition, model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port)
	defer stmt.Close()
	handleError("query failed: %v", err)
	affecteds, err := res.RowsAffected()
	handleError("could not get affected rows: %v", err)
	check := affecteds

	if check > 0 {
		return true
	}
	return false
}

// Function to extract values from the database depending on specific query parameters and bind them to a
func (serviceQueryList *ServiceQueryList) ServiceDefenitionFilter(serviceQueryForm ServiceQueryForm) {
	var queryHits []ServiceRegistryEntryOutput

	rows, err := db.Query("SELECT * FROM Services WHERE serviceDefinition LIKE ?", "%"+serviceQueryForm.ServiceDefinitionRequirement+"%")
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

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
		if validityCheck(service.EndOfValidity) {

			queryHits = append(queryHits, service)
		}

	}

	for i := 0; i < len(queryHits); i++ {

		queryHits[i].Interfaces = getInterfaceByID(queryHits[i].ID)
		queryHits[i].MetadataJava = getMetadataByID(queryHits[i].ID)
	}

	serviceQueryList.ServiceQueryData = queryHits
	serviceQueryList.UnfilteredHits = len(queryHits)

	return

}

// Function that binds services corresponding to the specific metadata variables in the form.
func (serviceQueryList *ServiceQueryList) metadataRequiermentFilter(serviceQueryForm ServiceQueryForm) {
	var metadataHits []ServiceRegistryEntryOutput
	for _, service := range serviceQueryList.ServiceQueryData {
		mdReqHit := false
		for serviceMdKey, serviceMdValue := range service.MetadataJava {
			for mdReqKey, mdReqValue := range serviceQueryForm.MetadataRequirementsJava {
				if serviceMdKey == mdReqKey || serviceMdValue == mdReqValue {
					mdReqHit = true

				}
			}
		}
		if mdReqHit {
			metadataHits = append(metadataHits, service)
		}
	}
	serviceQueryList.ServiceQueryData = metadataHits
	return
}

//Function to get a specific or all services depending id number in database, if id <= 0 then all services will be retured
func getServiceByID(id int64) []ServiceRegistryEntryOutput { //
	//var err error

	var serviceList []ServiceRegistryEntryOutput

	var rows *sql.Rows
	var err error
	if 0 <= id {
		rows, err = db.Query("SELECT * FROM Services WHERE id=?", id)
		defer rows.Close()
	} else {
		rows, err = db.Query("SELECT * FROM Services")
		defer rows.Close()
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
		if validityCheck(service.EndOfValidity) {

			serviceList = append(serviceList, service)
		}

	}
	for _, service := range serviceList {
		service.MetadataJava = getMetadataByID(service.ID)
		service.Interfaces = getInterfaceByID(service.ID)
	}
	return serviceList
}

//Function to get a specific or all interfaces corresponding to an id in the database.
func getInterfaceByID(id int) []Interface { //If id <= 0  then all interfaces will be retured

	var interfaces []Interface
	var rows *sql.Rows
	var err error
	if 0 <= id {
		rows, err = db.Query("SELECT serviceID, interfaceName, createdAt, updatedAt FROM Interfaces WHERE serviceID = ?", id)
		defer rows.Close()
	} else {
		rows, err = db.Query("SELECT serviceID, interfaceName, createdAt, updatedAt FROM Interfaces")
		defer rows.Close()
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

	return interfaces
}

//Function to get a specific or all metadata corresponding to an id in the database.
func getMetadataByID(id int) map[string]string { //If id <= 0 then all metadata will be retured

	metaData := make(map[string]string)
	var rows *sql.Rows

	rows, err := db.Query("SELECT key, value  FROM Metadata WHERE serviceID = ?", id)
	defer rows.Close()

	handleError("query failed: %v", err)
	for rows.Next() {
		var key string
		var value string

		err := rows.Scan(&key, &value)
		handleError("Encounterd an error while scaning database for metadata: %v", err)
		metaData[key] = value

	}

	return metaData
}

// Helper function to delete service for the validity check.
func deleteByID(id int) {

	stmt, err := db.Prepare("DELETE FROM Services WHERE id = ?")
	handleError("could not prepare statement: %v", err)
	_, err = stmt.Exec(id)
	defer stmt.Close()

	handleError("delete failed: %v", err)
}

// Function to remove/unregister outdated service from the database.
func cleanPastValidityDate() {
	fmt.Println("Preforming endOfValidity Cleaning")
	var pastValidityServices []int
	var id int
	var endOfValidity string

	rows, err := db.Query("SELECT id, endOfValidity FROM Services")
	defer rows.Close()
	if err != nil {

		panic(err.Error())
	}

	for rows.Next() {

		err := rows.Scan(&id, &endOfValidity)
		if err != nil {
			panic("query failed (endOfValidity cleaning): " + err.Error())
		}

		if !validityCheck(endOfValidity) {
			pastValidityServices = append(pastValidityServices, id)
		}
	}
	for _, v := range pastValidityServices {

		fmt.Printf("Removing service with ID: %d, the endOfValidity is not valid\n", v)
		deleteByID(v)
	}

}
func validityCheck(timeString string) bool {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {

		return false
	}
	return time.Now().Before(t)

}

// Function to set a ticker frequence for the time between validity cleaning
func startValidityTimer(minutes int) {
	cleanPastValidityDate()
	fmt.Printf("Clean delay set to: %d minutes\n", minutes)
	ticker := time.NewTicker(time.Duration(minutes) * time.Minute)
	for _ = range ticker.C {

		cleanPastValidityDate()
	}
}

// Function to get the array version of metadata
// (used when responding to requests with a MetadataGo containing entries.)
func convertToGoMetadata(metadata map[string]string) []string {
	var metadataGo = []string{}
	for key, value := range metadata {
		metadataGo = append(metadataGo, key+"."+value)
	}
	return metadataGo
}

// Function to get the map version of metadata
// (used when a requests with a MetadataGo contains entries.)
func convertToJavaMetadata(metadataGo []string) map[string]string {
	metadata := make(map[string]string)
	for _, key := range metadataGo {
		metadata[key] = key
	}
	return metadata
}
