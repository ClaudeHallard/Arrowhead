package ServiceRegistry

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

//var db *sql.DB

// OpenDatabase is functions to open database connectivity

func OpenDatabase() *sql.DB {
	var db *sql.DB
	var err error
	db, err = sql.Open("sqlite3", "file:ServiceRegistry/database/registryDB.db?cache=shared&_journal=WAL&_foreign_keys=on")
	if err != nil {
		panic(err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	return db
}

func handleError(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

//Query
func (model *ServiceQueryForm) Query() *ServiceQueryList {
	if len(model.MetadataRequirementsGo) == 0 {
		println("converting")
		model.MetadataRequirementsGo = convertToArrayFromStruct(model.MetadataRequirementsJava)
	}
	serviceQueryList := &ServiceQueryList{}
	//serviceQueryList.ServiceQueryData = getServiceByID(-1)
	serviceQueryList.ServiceDefenitionFilter(*model)

	if len(model.MetadataRequirementsGo) > 0 {

		serviceQueryList.metadataRequiermentFilter(*model)
	}
	for i := 0; i < len(serviceQueryList.ServiceQueryData); i++ {
		serviceQueryList.ServiceQueryData[i].MetadataJava = convertToStructFromArray(serviceQueryList.ServiceQueryData[i].MetadataGo)
	}

	return serviceQueryList
}

func (model *ServiceRegistryEntryInput) updateService() *ServiceRegistryEntryOutput {
	db := OpenDatabase()
	stmt, err := db.Prepare("UPDATE Services SET (systemName, address, port, authenticationInfo, endOfValidity, secure, version, updatedAt) = (?,?,?,?,?,?,?,?) WHERE serviceDefinition = ? AND serviceURI = ?")

	currentTime := time.Now().Format("2006-01-02T15:04:051Z")
	_, err = stmt.Exec(model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port, model.ProviderSystem.AuthenticationInfo, model.EndOfvalidity, model.Secure, model.Version, currentTime, model.ServiceDefinition, model.ServiceUri)
	if err != nil {
		println(err.Error())
		panic("Encounterd an error during registration while inserting service, (updateService)")
	}

	var cnt int64
	_ = db.QueryRow(`select id from Services where serviceURI = ?`, model.ServiceUri).Scan(&cnt)
	defer db.Close()

	return &getServiceByID(cnt)[0]
}

// check for serviceURI AND / OR serviceDefinition
func GetCountUniqueURI(servicedef string, serviceURI string) int {
	/*
		var cnt int
		db := OpenDatabase()
		//_ = db.QueryRow(`select count(*) from Services where serviceDefinition= ? AND serviceURI = ?`, servicedef, serviceURI).Scan(&cnt)
		_ = db.QueryRow(`select count(*) from Services where serviceURI = ?`, serviceURI).Scan(&cnt)

		defer db.Close()
		return cnt
	*/
	var cnt int
	db := OpenDatabase()
	row, err := db.Query(`select count(*) from Services where serviceDefinition= ? AND serviceURI = ?`, servicedef, serviceURI)
	defer row.Close()
	row.Scan(&cnt)
	if err != nil {
		panic("Encounterd an error during GetCountUniqueURI" + err.Error())
	}
	defer db.Close()
	return cnt

}

// check for serviceURI AND / OR serviceDefinition
func GetCountUpdate(servicedef string, serviceURI string) int {
	/*
			db := OpenDatabase()
		var cnt int
		//_ = db.QueryRow(`select count(*) from Services where serviceDefinition= ? AND serviceURI = ?`, servicedef, serviceURI).Scan(&cnt)
		_ = db.QueryRow(`select count(*) from Services where serviceURI = ? AND serviceDefinition = ?`, serviceURI, servicedef).Scan(&cnt)
		defer db.Close()
		return cnt

	*/
	db := OpenDatabase()
	var cnt int
	//_ = db.QueryRow(`select count(*) from Services where serviceDefinition= ? AND serviceURI = ?`, servicedef, serviceURI).Scan(&cnt)
	row, err := db.Query(`select count(*) from Services where serviceURI = ? AND serviceDefinition = ?`, serviceURI, servicedef)
	defer row.Close()
	row.Scan(&cnt)
	if err != nil {
		panic("Encounterd an error during GetCountUpdate" + err.Error())
	}
	defer db.Close()
	return cnt
}

// function to store the register input
func (model *ServiceRegistryEntryInput) Save() *ServiceRegistryEntryOutput {

	if !validityCheck(model.EndOfvalidity) {
		return nil
	}
	if len(model.MetadataGo) == 0 {
		println("converting")
		model.MetadataGo = convertToArrayFromStruct(model.MetadataJava)
	}

	// check if the service exist
	if GetCountUniqueURI(model.ServiceDefinition, model.ServiceUri) > 0 {
		if GetCountUpdate(model.ServiceDefinition, model.ServiceUri) > 0 {
			println("Register worked, service updated")
			return model.updateService() // update the service in the database

		}
		return nil

	} else {
		db := OpenDatabase()
		stmt, err := db.Prepare("INSERT INTO Services (serviceDefinition, systemName, address, port, authenticationInfo, serviceURI, endOfValidity, secure, version, createdAt, updatedAt) VALUES (?,?,?,?,?,?,?,?,?,?,?)")

		currentTime := time.Now().Format("2006-01-02T15:04:05.000Z")

		res, err := stmt.Exec(model.ServiceDefinition, model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port, model.ProviderSystem.AuthenticationInfo, model.ServiceUri, model.EndOfvalidity, model.Secure, model.Version, currentTime, currentTime)
		stmt.Close()
		if err != nil {
			println(err.Error())
			panic("Encounterd an error during registration while inserting service, (Save)")
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
			stmt.Close()

		}

		//loop through the interfaces array and add them to the table.
		for _, v := range model.Interfaces {
			stmt, err := db.Prepare("INSERT INTO Interfaces (serviceID, interfaceName, createdAt, updatedAt) VALUES (?,?,?,?)")
			handleError("could prepare statement: %v", err)
			_, err = stmt.Exec(lastId, v, currentTime, currentTime)
			stmt.Close()
		}

		handleError("failed to store: %v", err)
		println("Register worked")

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
		//println(len(model.MetadataGo))
		returnForm.MetadataJava = convertToStructFromArray(returnForm.MetadataGo)
		defer db.Close()
		return &returnForm

		//return &getServiceByID(lastId)[0]
	}
	//return nil
}

// Function to delete a service and the assosiated data
func (model *ServiceRegistryEntryInput) Delete() bool {
	db := OpenDatabase()
	stmt, err := db.Prepare("DELETE FROM Services WHERE serviceDefinition = ? AND systemName = ? AND  address = ? AND port = ?")
	handleError("could not prepare statement: %v", err)
	res, err := stmt.Exec(model.ServiceDefinition, model.ProviderSystem.SystemName, model.ProviderSystem.Address, model.ProviderSystem.Port)
	handleError("query failed: %v", err)
	affecteds, err := res.RowsAffected()
	handleError("could not get affected rows: %v", err)
	check := affecteds
	defer db.Close()
	if check > 0 {
		return true
	}
	return false
}
func (serviceQueryList *ServiceQueryList) ServiceDefenitionFilter(serviceQueryForm ServiceQueryForm) {
	var queryHits []ServiceRegistryEntryOutput
	db := OpenDatabase()
	rows, err := db.Query("SELECT * FROM Services WHERE serviceDefinition LIKE ?", "%"+serviceQueryForm.ServiceDefinitionRequirement+"%")
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
	defer rows.Close()
	defer db.Close()
	for i := 0; i < len(queryHits); i++ {
		queryHits[i].Interfaces = getInterfaceByID(int64(queryHits[i].ID))
		_, queryHits[i].MetadataGo = getMetadataByID(int64(queryHits[i].ID))
	}

	println(len(queryHits[0].MetadataGo))
	serviceQueryList.ServiceQueryData = queryHits
	serviceQueryList.UnfilteredHits = len(queryHits)

	return

}

func (serviceQueryList *ServiceQueryList) metadataRequiermentFilter(serviceQueryForm ServiceQueryForm) {
	var metadataHits []ServiceRegistryEntryOutput
	for _, service := range serviceQueryList.ServiceQueryData {
		mdReqHit := false
		for _, md := range service.MetadataGo {
			for _, mdReq := range serviceQueryForm.MetadataRequirementsGo {
				if strings.Contains(md, mdReq) {
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

//Function to get a specific or all services
func getServiceByID(id int64) []ServiceRegistryEntryOutput { //If id <= 0 then all services will be retured
	//var err error
	db := OpenDatabase()
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
		if validityCheck(service.EndOfValidity) {
			for i := 0; i < len(metaData); i++ {

				if service.ID == metaServiceID[i] {
					service.MetadataGo = append(service.MetadataGo, metaData[i])
				}
			}
			for _, v := range interfaces {
				if service.ID == v.ID {
					service.Interfaces = append(service.Interfaces, v)
				}
			}

			serviceList = append(serviceList, service)
		}

	}

	defer rows.Close()
	defer db.Close()
	return serviceList
}

//Function to get a specific or all interfaces
func getInterfaceByID(id int64) []Interface { //If id <= 0  then all interfaces will be retured
	db := OpenDatabase()
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
	defer db.Close()
	return interfaces
}

//Function to get a specific or all metadata
func getMetadataByID(id int64) ([]int, []string) { //If id <= 0 then all metadata will be retured
	var serviceID []int
	var metaData []string
	var rows *sql.Rows
	var err error
	db := OpenDatabase()
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
	defer db.Close()
	return serviceID, metaData
}
func deleteByID(id int) {
	db := OpenDatabase()
	stmt, err := db.Prepare("DELETE FROM Services WHERE id = ?")
	handleError("could not prepare statement: %v", err)
	_, err = stmt.Exec(id)

	handleError("delete failed: %v", err)
	defer db.Close()
}
func cleanPastValidityDate() {
	fmt.Println("Preforming endOfValidity Cleaning")
	var pastValidityServices []int
	var id int
	var endOfValidity string
	db := OpenDatabase()
	rows, err := db.Query("SELECT id, endOfValidity FROM Services")

	if err != nil {

		panic(err.Error())
	}

	for rows.Next() {

		err := rows.Scan(&id, &endOfValidity)
		handleError("query failed: %v", err)
		if !validityCheck(endOfValidity) {
			pastValidityServices = append(pastValidityServices, id)
		}
	}
	defer rows.Close()
	for _, v := range pastValidityServices {

		fmt.Printf("Removing service with ID: %d, the endOfValidity is not valid\n", v)
		deleteByID(v)
	}
	defer db.Close()

}
func validityCheck(timeString string) bool {
	t, err := time.Parse(time.RFC3339, timeString)
	if err != nil {

		return false
	}
	return time.Now().Before(t)

}

//ticker for the time between validity cleaning
func startValidityTimer(minutes int) {
	cleanPastValidityDate()
	fmt.Printf("Clean delay set to: %d minutes\n", minutes)
	ticker := time.NewTicker(time.Duration(minutes) * time.Minute)
	for _ = range ticker.C {

		cleanPastValidityDate()
	}
}
func convertToArrayFromStruct(metadataJava MetadataJava) []string {
	var stringArr []string
	if metadataJava.AdditionalProp1 != "" {
		stringArr = append(stringArr, metadataJava.AdditionalProp1)
	}
	if metadataJava.AdditionalProp2 != "" {
		stringArr = append(stringArr, metadataJava.AdditionalProp2)
	}
	if metadataJava.AdditionalProp3 != "" {
		stringArr = append(stringArr, metadataJava.AdditionalProp3)
	}
	return stringArr
}
func convertToStructFromArray(metadata []string) MetadataJava {
	metadataJava := MetadataJava{}

	if len(metadata) >= 1 {
		metadataJava.AdditionalProp1 = metadata[0]
	}
	if len(metadata) >= 2 {
		metadataJava.AdditionalProp2 = metadata[1]
	}
	if len(metadata) >= 3 {
		metadataJava.AdditionalProp3 = metadata[2]
	}
	return metadataJava

}
