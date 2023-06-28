package open_api

type OpenAPI struct {
	Info    Info     `json:"info"`
	Servers []Server `json:"servers"`
	Tags    []Tag    `json:"tags"`
	Paths   map[string]map[string]Method
	Components
}

type Info struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Version        string `json:"version"`
	TermsOfService string `json:"termsOfService"`
	Contact        struct {
		Name  string `json:"name"`
		Url   string `json:"url"`
		Email string `json:"email"`
	} `json:"contact"`
	License struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"license"`
}

type Server struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type Tag struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ExternalDocs struct {
		Description string `json:"description"`
		Url         string `json:"url"`
	} `json:"externalDocs,omitempty"`
}

type Method struct {
	Tags        []string `json:"tags"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	OperationId string   `json:"operationId"`
	Parameters  []Parameter
	Responses   []map[string]Response       `json:"responses"`
	Security    []map[string]SecurityScheme `json:"security"`
}

type Components struct {
	schemas string

	responses string

	Parameters []Parameter `json:"parameters"`

	examples string

	requestBodies string

	headers string

	SecuritySchemes []map[string]SecurityScheme `json:"securitySchemes"`

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
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	Description string    `json:"description,omitempty"`
	Format      string    `json:"format,omitempty"`
	Ref         string    `json:"$ref,omitempty"`
	Required    []string  `json:"required,omitempty"`
	Items       *Property `json:"items,omitempty"`
	File        *Property `json:"file,omitempty"`
	Enum        []string  `json:"enum,omitempty"`
}

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
}

type Header struct {
	Description string   `json:"description,omitempty"`
	Schema      Property `json:"schema,omitempty"`
}
