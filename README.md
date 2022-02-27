## Go's implementation of the Arrowhead framework

# Service Registry

##  __System Overview__

This System is a small demo version of the [Arrowhead Framework 4.4.0](https://github.com/eclipse-arrowhead/core-java-spring)'s service registry implemented in Golang, an SQLite based storage is used for this demo. The demo is both compatible with the original Framework implemented in Java (info listed in [additional doc](#additional-documentation)) and another Go-based implementation with two other demos, one [Orchestrator](https://github.com/Cordobro/D0020E)  and one [Application system](https://github.com/DavidS1998/D00020E) part. 

<br>

This project is a part of the Course D0020E at Luleå University of Technology 2021/2022

<br>

## Design structures and UML schematics:
#### The Basic design for this service registry can be found in the [UML directory](https://github.com/ClaudeHallard/Arrowhead/tree/main/UML) , it contains all diagrams and figures stated below as well as the raw UML format made in [Eclipse Papyrus™](https://www.eclipse.org/papyrus/).
- [System Overview](https://raw.githubusercontent.com/ClaudeHallard/Arrowhead/main/UML/images/Overview.png) (Relation view for the service registry and other API's in the framework)

- [Use Case Diagram](https://raw.githubusercontent.com/ClaudeHallard/Arrowhead/main/UML/images/UseCaseDiagram.png)

- [Class Diagram](https://raw.githubusercontent.com/ClaudeHallard/Arrowhead/main/UML/images/class_fileStructureDiagram.png)

- [Database Schema](https://raw.githubusercontent.com/ClaudeHallard/Arrowhead/main/UML/images/DatabaseSchema.png)

- [Activity Diagram](https://raw.githubusercontent.com/ClaudeHallard/Arrowhead/main/UML/images/ProviderActivityDiagram.png) (Provider client)

- [Activity Diagram](https://raw.githubusercontent.com/ClaudeHallard/Arrowhead/main/UML/images/ConsumerActivityDiagram.png) (Consumer client)

- [State Diagram](https://raw.githubusercontent.com/ClaudeHallard/Arrowhead/main/UML/images/StateDiagram.png)


<br>

## Compile and run instructions

#### __Windows users:__ run the "```serviceRegistry.exe```"-file or run and compile with ```go build && ./serviceRegistry``` in the "src" directory, it can also be compiled without building with ```go run starter.go```.
##### OBS!, Make sure to have a c-compiler installed in your %PATH to be able to use the imported SQLite extenstion used in this project.
<br>

#### __MAC/unix users:__ compile and run with ```go run starter.go``` in the "src" directory or ```go build && ./serviceRegistry``` if the system supports the creation of executable files.

<br>

__Note1:__ depending on the current go setup and version you may need to delete the .mod file and run  ```go mod init``` and ```go mod tidy``` to compile the program.

__Note2:__ The ```TestClient.go``` folder contains example payloads and test that can be used to to run diagnostics on this serviceRegistry, it contains methods for testing all services provided in this demo, it can be compiled with ```go run *.go``` in the test directory.


### __Database clarifications__
This implementation is using a local based SQLite database stored in the file: ```registryDB.db``` in the Database directory. No database management system is required to use the Service Registry but is recommended for eventual testing and debug checks. In the developement of this demo, ["__SQLitestudio©__"](https://sqlitestudio.pl/) and ["__DBBrowser__"](https://sqlitebrowser.org/) has been used.


<br>

## Additional documentation and information

### __Imported repos__
The project is importing two other Github repositories for its implementaion, the ["github.com/hariadivicky/nano"](https://github.com/hariadivicky/nano) repo for an easy use of RESTful methods and for storing JSON data into the database. The interface between golang and SQL is imported with the ["github.com/mattn/go-sqlite3"](https://github.com/mattn/go-sqlite3) repo.

### __GoDoc - documentaion__
The file and folder structure supports GoDoc generation for this project repo, setup guide can be found [HERE](https://pkg.go.dev/golang.org/x/tools/cmd/godoc)

### __Support for the Original Java Implementaion__



Payloads used in the Java [Arrowhead Framework](https://github.com/eclipse-arrowhead/core-java-spring) can also be used in this Golang version, the strucs in [```JSONForms.go```](https://github.com/ClaudeHallard/Arrowhead/blob/main/src/ServiceRegistry/JSONForms.go) contain specific fields for "MetaDataJava" and "MetaDataGo" and the functionalites for extracting data for each of the corresponding payloads. See the payload struct listed for each method for more details.

### __Test functions__

Comparation between the [Arrowhead Framework](https://github.com/eclipse-arrowhead/core-java-spring) and this demo has been implemented in this project, the tests can be found in the ```TestClient.go``` -file where also payload tests can be ran. 

__Note:__ The Java program is not provided in this repository

<br>
<br>

## Service methods Overview
# 
The service Registry supports four operands for client systems through HTTP requests, 

The __echo__ method

The __register__ method 

The __query__ method 

The __unregister__ method 

__Note:__  There are other functionalies and methods in the original Arrohead Framework's Service Registry that is not implemented in this demo.


The base URL for the requests: `http://<host>:<port>/serviceregistry` (for this demo the localhost is set to port 4245)

### Service description<br />

| Function | URL subpath | Method | Input | Output |
| -------- | ----------- | ------ | ----- | ------ |
| Echo    | /echo       | GET    | -     | OK     |
| Query   | /query      | POST   | ServiceQueryForm | ServiceQueryList |
| Register | /register   | POST   | ServiceRegistryEntryInput | ServiceRegistryEntryOutput |
| Unregister | /unregister | DELETE | Address, Port, Service Definition, System Name in query parameters| OK |


# Echo 
```
GET /serviceregistry/echo
```

Returns a confirmation message with the purpose of testing if the Service Registry is online

# Register
```
POST /serviceregistry/register
```

Registers a service from a provider.

__ServiceRegistryEntryInput__ is the input (Go struct with json format)
```
type ServiceRegistryEntryInput struct {
	ServiceDefinition string         `json:"serviceDefinition"`
	ProviderSystem    ProviderSystem `json:"providerSystem"`
	ServiceUri        string         `json:"serviceUri"`
	EndOfvalidity     string         `json:"endOfValidity"`
	Secure            string         `json:"secure"`
	MetadataGo        []string       `json:"metadataGo"`
	MetadataJava      MetadataJava   `json:"metadata"`
	Version           int            `json:"version"`
	Interfaces        []string       `json:"interfaces"`
}
type ProviderSystem struct {
	SystemName         string `json:"systemName"`
	Address            string `json:"address"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
}
```
## Input clarifiction:
| Field | Description | Mandatory |
| ----- | ----------- | --------- |
| `serviceDefinition` | Service Definition | yes |
| `providerSystem` | Provider System | yes |
| `serviceUri` |  URI of the service | yes |
| `endOfValidity` | TimeStamp, must be valid and in correct format  | yes |
| `secure` | Security info | no |
| `metadataGo` | Metadata for this go demo | no |
| `metadataJava` | Metadata for the Java version| no |
| `version` | Version of the Service | no |
| `interfaces` | List of the interfaces the Service supports | no |

Returns a __ServiceRegistryEntryOutput (Go struct with json format)__

```
type ServiceRegistryEntryOutput struct {
	ID                int               `json:"id"`
	ServiceDefinition ServiceDefinition `json:"serviceDefinition"`
	Provider          Provider          `json:"provider"`
	ServiceUri        string            `json:"serviceUri"`
	EndOfValidity     string            `json:"endOfValidity"`
	Secure            string            `json:"secure"`
	MetadataGo        []string          `json:"metadataGo"`
	MetadataJava      MetadataJava      `json:"metadata"`
	Version           int               `json:"version"`
	Interfaces        []Interface       `json:"interfaces"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
}

type ServiceDefinition struct {
	ID                int    `json:"id"`
	ServiceDefinition string `json:"serviceDefinition"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}
type Provider struct {
	ID                 int    `json:"id"`
	SystemName         string `json:"systemName"`
	Address            string `json:"adress"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type Interface struct {
	ID            int    `json:"id"`
	InterfaceName string `json:"interfaceName"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}
```
## Output clarifiction:
| Field | Description |
| ----- | ----------- |
| `id` | ID of the ServiceRegistryEntry (internal for the databse) |
| `serviceDefinition` | Service Definition |
| `provider` | Provider System |
| `serviceUri` |  URI of the Service |
| `endOfValidity` | Service is available until this UTC timestamp |
| `secure` | Security info |
| `metadataGo` | Metadata for this go demo |
| `metadataJava` | Metadata for the Java version|
| `version` | Version of the Service |
| `interfaces` | List of the interfaces the Service supports |
| `createdAt` | Creation date of the entry |
| `updatedAt` | When the entry was last updated (often same as creation) |


# Query
```
POST /serviceregistry/query
```

Returns ServiceQueryList for the specific service. If no service was found, an error message is sent to the client


__ServiceQueryForm__ is the input (Go struct with json format)
```
type ServiceQueryForm struct {
	ServiceDefinitionRequirement string       `json:"serviceDefinitionRequirement"`
	InterfaceRequirements        []string     `json:"interfaceRequirements"`
	SecurityRequirements         []string     `json:"securityRequirements"`
	MetadataRequirementsGo       []string     `json:"metadataRequirementsGo"`
	MetadataRequirementsJava     MetadataJava `json:"metadataRequirements"`
	VersionRequirements          int          `json:"versionRequirement"`
	MaxVersionRequirements       int          `json:"maxVersionRequirement"`
	MinVersionRequirements       int          `json:"minVersionRequirement"`
	PingProviders                bool         `json:"pingProviders"`
}
```
## Input clarifiction:
| Field | Description | Mandatory |
| ----- | ----------- | --------- |
| `serviceDefinitionRequirement` | Name of the required Service Definition | yes |
| `interfaceRequirements` | List of required interfaces | no |
| `securityRequirements` | List of required security settings | no |
| `metadataRequirementsGo` | Value of required metadata | no |
| `metadataRequirementsJava` | Key value pairs of required metadata | no |
| `versionRequirement` | Required version number | no |
| `maxVersionRequirement` | Maximum version requirement | no |
| `minVersionRequirement` | Minimum version requirement | no |
| `pingProviders` | Return only available providers | no |


Returns a __ServiceQueryList (Go struct with json format)__
```
type ServiceQueryList struct {
	ServiceQueryData []ServiceRegistryEntryOutput `json:"serviceQueryData"`
	UnfilteredHits   int                          `json:"unfilteredHits"`
}

type ServiceRegistryEntryOutput struct {
	ID                int               `json:"id"`
	ServiceDefinition ServiceDefinition `json:"serviceDefinition"`
	Provider          Provider          `json:"provider"`
	ServiceUri        string            `json:"serviceUri"`
	EndOfValidity     string            `json:"endOfValidity"`
	Secure            string            `json:"secure"`
	MetadataGo        []string          `json:"metadataGo"`
	MetadataJava      MetadataJava      `json:"metadata"`
	Version           int               `json:"version"`
	Interfaces        []Interface       `json:"interfaces"` 
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
}

type ServiceDefinition struct {
	ID                int    `json:"id"`
	ServiceDefinition string `json:"serviceDefinition"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type Provider struct {
	ID                 int    `json:"id"`
	SystemName         string `json:"systemName"`
	Address            string `json:"adress"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type Interface struct {
	ID            int    `json:"id"`
	InterfaceName string `json:"interfaceName"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}
```
## Output clarifiction:
| Field | Description |
| ----- | ----------- |
| `serviceQueryData` | The array of objects containing the data |
| `id` | ID of the entry, used by the Orchestrator |
| `serviceDefinition` | Service Definition |
| `provider` | Provider System |
| `serviceUri` | URI of the Service |
| `endOfValidity` | Service is available until this UTC timestamp. |
| `secure` | Security info |
| `metadataGo` | Metadata payload for Go|
| `metadataJava` | Metadata payload for Java|
| `version` | Version of the Service |
| `interfaces` | List of interfaces the Service supports |
| `createdAt` | Creation date of the entry |
| `updatedAt` | When the entry was last updated |
| `unfilteredHits` | Number of hits based on service definition without filters |


# Unregister 
```
DELETE /serviceregistry/unregister
```

Removes a registered service, returns an "OK" if service was removed and an "FAILED" if it falied or 

Pure URI provided data is the input

Example format where `"X"`, `"Y"`, `"Z"` ,`"Q"` and `"W"`is provided data:
`
http://<host>:<port>/serviceregistry/unregister?address="X"&port="Y"&service_definition="Z"&service_uri="Q"&system_name="W"
`
## Input clarifiction:

| Field | Description | Mandatory |
| ----- | ----------- | --------- |
| `address` | Service address | yes |
| `port` | Port number | yes |
| `service_definition` | Service Definition | yes |
| `service_uri` |  URI of the service | yes |
| `system_name` |  Name of the system | yes |


Note: Inorder to unregistrer a service,  `adress`, `port`, `serviceDefinition` , `service_uri` &  `systemName` must be all be provided and valid

