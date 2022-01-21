# Go's implementation of the Arrowhead framework

## Service Registry

##  Overview

This System provides the database, which stores information related to the currently actively offered Services within the Local Cloud.


<a name="serviceregistry_usecases" />

## Service methods Overview

The __echo__ method

The __register__ method 

The __query__ method 

The __unregister__ method 

There are other functionalies of the arrowhead framework that is not implemented in this demo of the Service Registry.

## Endpoints

The Service Registry offers three types of endpoints. Client, Management and Private.

Swagger API documentation is available on: `https://<host>:<port>` <br />
The base URL for the requests: `http://<host>:<port>/serviceregistry`

### Client endpoint description<br />

| Function | URL subpath | Method | Input | Output |
| -------- | ----------- | ------ | ----- | ------ |
| Echo    | /echo       | GET    | -     | OK     |
| Query   | /query      | POST   | ServiceQueryForm | ServiceQueryList |
| Register | /register   | POST   | ServiceRegistryEntry | ServiceRegistryEntry |
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

Registers a service.

__ServiceRegistryEntryInput__ is the input (Go struct with json format)
```json
type ServiceRegistryEntryInput struct {
	ServiceDefinition string         `json:"serviceDefinition"`
	ProviderSystem    ProviderSystem `json:"providerSystem"`
	ServiceUri        string         `json:"serviceUri"`
	EndOfvalidity     string         `json:"endOfValidity"`
	Secure            string         `json:"NOT_SECURE"`
	Metadata          Metadata       `json:"metadata"`
	Version           int            `json:"version"`
	Interfaces        []string       `json:"interfaces"`
}
type ProviderSystem struct {
	SystemName         string `json:"systemName"`
	Address            string `json:"adress"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
}

type Metadata struct {
	AdditionalProp1 string `json:"additionalProp1"`
	AdditionalProp2 string `json:"additionalProp2"`
	AdditionalProp3 string `json:"additionalProp3"`
}
```
## Input clarifiction:
| Field | Description | Mandatory |
| ----- | ----------- | --------- |
| `serviceDefinition` | Service Definition | yes |
| `providerSystem` | Provider System | yes |
| `serviceUri` |  URI of the service | yes |
| `endOfValidity` | Service is available until this UTC timestamp | yes |
| `secure` | Security info | yes |
| `metadata` | Metadata | yes |
| `version` | Version of the Service | yes |
| `interfaces` | List of the interfaces the Service supports | yes |

Returns a __ServiceRegistryEntryOutput (Go struct with json format)__

```json
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
| `id` | ID of the ServiceRegistryEntry |
| `serviceDefinition` | Service Definition |
| `provider` | Provider System |
| `serviceUri` |  URI of the Service |
| `endOfValidity` | Service is available until this UTC timestamp |
| `secure` | Security info |
| `metadata` | Metadata |
| `version` | Version of the Service |
| `interfaces` | List of the interfaces the Service supports |
| `createdAt` | Creation date of the entry |
| `updatedAt` | When the entry was last updated |


# Query
```
POST /serviceregistry/query
```

Returns ServiceQueryList 


__ServiceQueryForm__ is the input (Go struct with json format)
```json
type ServiceQueryForm struct {
	ServiceDefinitionRequirement string   `json:"serviceDefinitionRequirement"`
	InterfaceRequirements        []string `json:"interfaceRequirements"`
	SecurityReRequirements       []string `json:"securityRequirements"`
	MetadataRequirements         Metadata `json:"metadataRequirements "`
	VersionRequirements          int      `json:"versionRequirements"`
	MaxVersionRequirements       int      `json:"maxVersionRequirements"`
	MinVersionRRequirements      int      `json:"minVersionRRequirements"`
	PingProviders                bool     `json:"pingProviders"`
}

type Metadata struct {
	AdditionalProp1 string `json:"additionalProp1"`
	AdditionalProp2 string `json:"additionalProp2"`
	AdditionalProp3 string `json:"additionalProp3"`
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
```json
type ServiceQueryList struct {
	ServiceQueryData []ServiceRegistryEntryOutput `json:"serviceQueryData"`
	UnfilteredHits   int                          `json:"unfilteredHits"`
}

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

type Metadata struct {
	AdditionalProp1 string `json:"additionalProp1"`
	AdditionalProp2 string `json:"additionalProp2"`
	AdditionalProp3 string `json:"additionalProp3"`
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

Query params:

| Field | Description | Mandatory |
| ----- | ----------- | --------- |
| `service_definition` | Name of the service to be removed | yes |
| `system_name` | Name of the Provider | no |
| `address` | Address of the Provider | no |
| `port` | Port of the Provider | no |


