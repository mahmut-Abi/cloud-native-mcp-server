package openapi

// Types define the OpenAPI specification structures

type OpenAPISpec struct {
	OpenAPI    string                `json:"openapi"`
	Info       Info                  `json:"info"`
	Servers    []Server              `json:"servers"`
	Paths      map[string]PathItem   `json:"paths"`
	Components Components            `json:"components"`
	Tags       []Tag                 `json:"tags"`
	Security   []map[string][]string `json:"security,omitempty"`
}

type Info struct {
	Title       string                  `json:"title"`
	Version     string                  `json:"version"`
	Description string                  `json:"description"`
	Contact     *map[string]interface{} `json:"contact,omitempty"`
	License     *map[string]interface{} `json:"license,omitempty"`
}

type Server struct {
	URL         string                            `json:"url"`
	Description string                            `json:"description"`
	Variables   map[string]map[string]interface{} `json:"variables,omitempty"`
}

type PathItem struct {
	Post *Operation `json:"post,omitempty"`
	Get  *Operation `json:"get,omitempty"`
}

type Operation struct {
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	OperationID string                `json:"operationId"`
	Tags        []string              `json:"tags"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string `json:"security,omitempty"`
}

type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Schema      interface{} `json:"schema"`
}

type RequestBody struct {
	Required bool                 `json:"required"`
	Content  map[string]MediaType `json:"content"`
}

type MediaType struct {
	Schema  interface{} `json:"schema"`
	Example interface{} `json:"example,omitempty"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type Components struct {
	Schemas         map[string]interface{}    `json:"schemas,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

type SecurityScheme struct {
	Type         string `json:"type"`
	Description  string `json:"description"`
	Name         string `json:"name,omitempty"`
	In           string `json:"in,omitempty"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
