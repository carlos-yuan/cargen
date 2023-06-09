package openapi

type OpenAPI struct {
	Openapi    string                       `json:"openapi,omitempty"`
	Info       Info                         `json:"info,omitempty"`
	Servers    []Server                     `json:"servers,omitempty"`
	Tags       []Tag                        `json:"tags,omitempty"`
	Paths      map[string]map[string]Method `json:"paths,omitempty"`
	Components Components                   `json:"components,omitempty"`
}

type Info struct {
	Title          string   `json:"title,omitempty"`
	Description    string   `json:"description,omitempty"`
	Version        string   `json:"version,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	Url   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

var DefaultInfo = OpenAPI{
	Openapi: "3.0.0",
	Components: Components{
		Schemas:         make(map[string]Property),
		Responses:       make(map[string]Response),
		Headers:         make(map[string]Header),
		RequestBodies:   make(map[string]RequestBody),
		SecuritySchemes: make(map[string]SecurityScheme),
	},
	Paths: make(map[string]map[string]Method),
}

type Server struct {
	Url         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

type Tag struct {
	Name         string        `json:"name,omitempty"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
}

type Method struct {
	Tags        []string                    `json:"tags,omitempty"`
	Summary     string                      `json:"summary,omitempty"`
	Description string                      `json:"description,omitempty"`
	OperationId string                      `json:"operationId,omitempty"`
	Parameters  []Parameter                 `json:"parameters,omitempty"`
	RequestBody RequestBody                 `json:"requestBody,omitempty"`
	Responses   map[string]Response         `json:"responses,omitempty"`
	Security    []map[string]SecurityScheme `json:"security,omitempty"`
	api         *OpenAPI
}

type Components struct {
	Schemas map[string]Property `json:"schemas,omitempty"`

	Responses map[string]Response `json:"responses,omitempty"`

	Parameters []Parameter `json:"parameters,omitempty"`

	RequestBodies map[string]RequestBody `json:"requestBodies,omitempty"`

	Headers map[string]Header `json:"headers,omitempty"`

	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

func (c Components) GetSchemasName() string {
	return "#/components/schemas/"
}

func (c Components) GetResponsesName() string {
	return "#/components/responses/"
}

func (c Components) GetParametersName() string {
	return "#/components/parameters/"
}

func (c Components) GetRequestBodiesName() string {
	return "#/components/requestBodies/"
}

func (c Components) GetHeadersName() string {
	return "#/components/headers/"
}

func (c Components) GetSecuritySchemesName() string {
	return "#/components/securitySchemes/"
}

type SecurityScheme struct {
	Type        string            `json:"type,omitempty"`
	Description string            `json:"description,omitempty"`
	Flows       []map[string]Flow `json:"flows,omitempty"`
}

type Flow struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	Scope            map[string]string `json:"scope,omitempty"`
}

type Property struct {
	Name        string              `json:"-"`
	Type        string              `json:"type,omitempty"`
	Description string              `json:"description,omitempty"`
	Format      string              `json:"format,omitempty"`
	Ref         string              `json:"$ref,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Items       *Property           `json:"items,omitempty"` //数组
	File        *Property           `json:"file,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
	Enum        []string            `json:"enum,omitempty"`
}

const (
	PropertyTypeObject = "object"
	PropertyTypeArray  = "array"
)

type Parameter struct {
	Name            string   `json:"name,omitempty"`
	In              string   `json:"in,omitempty"`
	Description     string   `json:"description,omitempty"`
	Required        bool     `json:"required,omitempty"`
	AllowEmptyValue bool     `json:"allowEmptyValue,omitempty"`
	Deprecated      bool     `json:"deprecated,omitempty"`
	Schema          Property `json:"schema,omitempty"`
	Style           string   `json:"style,omitempty"`
}

type RequestBody struct {
	Description string             `json:"description,omitempty"`
	Required    bool               `json:"required,omitempty"`
	Content     map[string]Content `json:"content,omitempty"`
	Ref         string             `json:"$ref,omitempty"`
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
