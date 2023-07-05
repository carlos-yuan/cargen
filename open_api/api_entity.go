package openapi

type OpenAPI struct {
	Openapi    string                       `json:"openapi"`
	Info       Info                         `json:"info,omitempty"`
	Servers    []Server                     `json:"servers,omitempty"`
	Tags       []Tag                        `json:"tags,omitempty"`
	Paths      map[string]map[string]Method `json:"paths,omitempty"`
	Components Components                   `json:"components,omitempty"`
}

type Info struct {
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	Version        string `json:"version,omitempty"`
	TermsOfService string `json:"termsOfService,omitempty"`
	Contact        struct {
		Name  string `json:"name"`
		Url   string `json:"url"`
		Email string `json:"email"`
	} `json:"contact"`
	License struct {
		Name string `json:"name,omitempty"`
		Url  string `json:"url,omitempty"`
	} `json:"license,omitempty"`
}

var DefaultInfo = OpenAPI{
	Openapi: "3.0.0",
}

type Server struct {
	Url         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

type Tag struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	ExternalDocs struct {
		Description string `json:"description,omitempty"`
		Url         string `json:"url,omitempty"`
	} `json:"externalDocs,omitempty"`
}

type Method struct {
	Tags        []string                    `json:"tags"`
	Summary     string                      `json:"summary"`
	Description string                      `json:"description"`
	OperationId string                      `json:"operationId"`
	Parameters  []Parameter                 `json:"parameters"`
	RequestBody RequestBody                 `json:"requestBody"`
	Responses   map[string]Response         `json:"responses"`
	Security    []map[string]SecurityScheme `json:"security"`
}

type Components struct {
	Schemas map[string]Property `json:"schemas"`

	Responses map[string]Response `json:"responses"`

	Parameters []Parameter `json:"parameters"`

	examples string

	RequestBodies map[string]RequestBody `json:"requestBodies"`

	Headers map[string]Header `json:"headers"`

	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes"`

	links string

	callbacks string
}

type SecurityScheme struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Flows       []map[string]Flow
}

type Flow struct {
	AuthorizationUrl string `json:"authorizationUrl,omitempty"`
	TokenUrl         string `json:"tokenUrl,omitempty"`
	Scope            map[string]string
}

type Property struct {
	Name        string              `json:"name,omitempty"`
	Type        string              `json:"type,omitempty"`
	Description string              `json:"description,omitempty"`
	Format      string              `json:"format,omitempty"`
	Ref         string              `json:"$ref,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Items       *Property           `json:"items,omitempty"`
	File        *Property           `json:"file,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
	Enum        []string            `json:"enum,omitempty"`
}

const (
	PropertyTypeObject = "object"
)

type Parameter struct {
	Name            string   `json:"name,omitempty"`
	In              string   `json:"in,omitempty"`
	Description     string   `json:"description,omitempty"`
	Required        bool     `json:"required,omitempty"`
	AllowEmptyValue bool     `json:"allowEmptyValue"`
	Deprecated      bool     `json:"deprecated"`
	Schema          Property `json:"schema,omitempty"`
	Style           string   `json:"style,omitempty"`
}

type RequestBody struct {
	Description string             `json:"description,omitempty"`
	Required    bool               `json:"required,omitempty"`
	Content     map[string]Content `json:"content,omitempty"`
}

type Content struct {
	Type   string   `json:"type,omitempty"`
	Schema Property `json:"schema,omitempty"`
}

type Response struct {
	Description string             `json:"description,omitempty"`
	Content     map[string]Content `json:"content,omitempty"`
	Headers     map[string]Header  `json:"headers,omitempty"`
}

type Header struct {
	Description string   `json:"description,omitempty"`
	Schema      Property `json:"schema,omitempty"`
}
