## Go's implementation of the Arrowhead framework

# Service Registry

##  __System Overview__

This System is a small demo version of the Arrowhead Framework 4.4.0's service registry implemented in GOLANG, an SQLite based storage is used for this demo. Design structures and schematics are provided in the UML directory. 


This project is a part of the Course D0020E at Luleå University of Technology
<br>

## Compile and run instructions

### __Windows users:__ run the "```serviceRegistry.exe```"-file or run 
### and compile with ```go build && ./serviceRegistry``` in the main directory
<br>

### __MAC/unix users:__ compile and run with ```go run *.go``` in the main directory

<br>

__Note:__ The ```clientTest.go``` file can be used to test and run diagnostics on this serviceRegistry demo, it contains different test methods for all services provided in this demo, it can be compiled with ```go run clientTest.go```
# 
## Service methods Overview

The __echo__ method

The __register__ method 

The __query__ method 

The __unregister__ method 

There are other functionalies of the arrowhead framework that is not implemented in this demo of the Service Registry.


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
	Metadata          []string       `json:"metadata"`
	Version           int            `json:"version"`
	Interfaces        []string       `json:"interfaces"`
}

type ProviderSystem struct {
	SystemName         string `json:"systemName"`
	Address            string `json:"adress"`
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
| `endOfValidity` | TimeStamp, overdue services will be unregistred  | no |
| `secure` | Security info | no |
| `metadata` | Metadata | no |
| `version` | Version of the Service | no |
| `interfaces` | List of the interfaces the Service supports | no |

Returns a __ServiceRegistryEntryOutput (Go struct with json format)__

```
type ServiceRegistryEntryOutput struct {
	ID                int               `json:"id"`
	ServiceDefinition ServiceDefinition `json:"serviceDefinition"`
	Provider          Provider          `json:"id"`
	ServiceUri        string            `json:"serviceUri"`
	EndOfvalidity     string            `json:"endOfValidity"`
	Secure            string            `json:"NOT_SECURE"`
	Metadata          Metadata          `json:"metadata"`
	Version           int               `json:"version"`
	interfaces        []Interface       `json:"interfaces"` 
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
| `metadata` | Metadata |
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
	ServiceDefinitionRequirement string   `json:"serviceDefinitionRequirement"`
	InterfaceRequirements        []string `json:"interfaceRequirements"`
	SecurityReRequirements       []string `json:"securityRequirements"`
	MetadataRequirements         []string `json:"metadataRequirements"`
	VersionRequirements          int      `json:"versionRequirement"`
	MaxVersionRequirements       int      `json:"maxVersionRequirement"`
	MinVersionRequirements       int      `json:"minVersionRequirement"`
	PingProviders                bool     `json:"pingProviders"`
}
```
## Input clarifiction:
| Field | Description | Mandatory |
| ----- | ----------- | --------- |
| `serviceDefinitionRequirement` | Name of the required Service Definition | yes |
| `interfaceRequirements` | List of required interfaces | no |
| `securityRequirements` | List of required security settings | no |
| `metadataRequirements` | Key value pairs of required metadata | no |
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
	Metadata          []string          `json:"metadata"`
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
| `metadata` | Metadata |
| `version` | Version of the Service |
| `interfaces` | List of interfaces the Service supports |
| `createdAt` | Creation date of the entry |
| `updatedAt` | When the entry was last updated |
| `unfilteredHits` | Number of hits based on service definition without filters |


# Unregister 
```
DELETE /serviceregistry/unregister
```

Removes a registered service, returns an OK if service was removed and an Error if failed

__ServiceRegistryEntryInput__ is the input (Go struct with json format)
```
type ServiceRegistryEntryInput struct {
	ServiceDefinition string         `json:"serviceDefinition"`
	ProviderSystem    ProviderSystem `json:"providerSystem"`
	ServiceUri        string         `json:"serviceUri"`
	EndOfvalidity     string         `json:"endOfValidity"`
	Secure            string         `json:"secure"`
	Metadata          []string       `json:"metadata"`
	Version           int            `json:"version"`
	Interfaces        []string       `json:"interfaces"`
}

type ProviderSystem struct {
	SystemName         string `json:"systemName"`
	Address            string `json:"adress"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
}

```
## Input clarifiction:

| Field | Description | Mandatory |
| ----- | ----------- | --------- |
| `serviceDefinition` | Service Definition | yes |
| `providerSystem` | Provider System | yes |
| `serviceUri` |  URI of the service | no |
| `endOfValidity` | TimeStamp, overdue services will be unregistred  | no |
| `secure` | Security info | no |
| `metadata` | Metadata | no |
| `version` | Version of the Service | no |
| `interfaces` | List of the interfaces the Service supports | no |

Note: (inorder to unregistrer a service `serviceDefinition` , `systemName` , `adress` & `port` must be provided and valid)

