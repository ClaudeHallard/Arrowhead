package main

//Shared Structs

//Metadata scruct, used by
//ServiceRegistryEntryInput, ServiceRegistryEntryOutput,
// ServiceQueryForm and ServiceQueryList.
type MetadataJava struct {
	AdditionalProp1 string `json:"additionalProp1,omitempty"`
	AdditionalProp2 string `json:"additionalProp2,omitempty"`
	AdditionalProp3 string `json:"additionalProp3,omitempty"`
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
	Address            string `json:"address"`
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

// ServiceRegistryEntry Input Version
type ServiceRegistryEntryInput struct {
	ServiceDefinition string         `json:"serviceDefinition"`
	ProviderSystem    ProviderSystem `json:"providerSystem"`
	ServiceUri        string         `json:"serviceUri"`
	EndOfvalidity     string         `json:"endOfValidity"`
	Secure            string         `json:"secure"`
	MetadataJava      MetadataJava   `json:"metadata"`
	MetadataGo        []string       `json:"metadataGo"`
	Version           int            `json:"version"`
	Interfaces        []string       `json:"interfaces"`
}
type ProviderSystem struct {
	SystemName         string `json:"systemName"`
	Address            string `json:"address"`
	Port               int    `json:"port"`
	AuthenticationInfo string `json:"authenticationInfo"`
}

// ServiceRegistryEntry Output Version
type ServiceRegistryEntryOutput struct {
	ID                int               `json:"id"`
	ServiceDefinition ServiceDefinition `json:"serviceDefinition"`
	Provider          Provider          `json:"provider"`
	ServiceUri        string            `json:"serviceUri"`
	EndOfValidity     string            `json:"endOfValidity"`
	Secure            string            `json:"secure"`
	MetadataJava      MetadataJava      `json:"metadata"`
	MetadataGo        []string          `json:"metadataGo"`
	Version           int               `json:"version"`
	Interfaces        []Interface       `json:"interfaces"` //I think this is how it's implemented -Ivar (should result in an array of interfaces.)
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
}

//ServiceQueryForm
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

// Test  func.
func TestFunc() {

}

// ServiceQueryList 1
type ServiceQueryList struct {
	ServiceQueryData []ServiceRegistryEntryOutput `json:"serviceQueryData"`
	UnfilteredHits   int                          `json:"unfilteredHits"`
}
