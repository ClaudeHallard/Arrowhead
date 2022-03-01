/*
**********************************************
@authors:
Ivar Grunaeu, ivagru-9@student.ltu.se
Jean-Claude Hallard, jeahal-8@student.ltu.se
Marcus Paulsson, marpau-8@student.ltu.se
Pontus Schünemann, ponsch-9@student.ltu.se
Fabian Widell, fabwid-9@student.ltu.se

**********************************************
JSONForms.go contains all the shared structs that contains all the data variables for the json forms */

package ServiceRegistry

/*Metadata scruct for the json field containing java variables, used by
ServiceRegistryEntryInput, ServiceRegistryEntryOutput,
ServiceQueryForm and ServiceQueryList. */
type MetadataJava struct {
	AdditionalProp1 string `json:"additionalProp1"`
	AdditionalProp2 string `json:"additionalProp2"`
	AdditionalProp3 string `json:"additionalProp3"`
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

// ServiceRegistryEntry Output Version
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

//ServiceQueryList contains a list of returned services
type ServiceQueryList struct {
	ServiceQueryData []ServiceRegistryEntryOutput `json:"serviceQueryData"`
	UnfilteredHits   int                          `json:"unfilteredHits"`
}
