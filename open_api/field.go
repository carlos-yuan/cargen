package openapi

import (
	"github.com/carlos-yuan/cargen/util"
	"reflect"
	"strings"
)

type Field struct {
	Name      string `json:"name"`      //名称
	Type      string `json:"type"`      //类型
	Tag       string `json:"tag"`       //tag
	In        string `json:"in"`        //存在类型 path 路径 json json对象内
	ParamName string `json:"paramName"` //参数名
	Validate  string `json:"validate"`  //验证
	Pkg       string `json:"pkg"`       //包名 类型为结构体时
	PkgPath   string `json:"pkgPath"`   //包路径 类型为结构体时
	Comment   string `json:"comment"`   //注释
	Array     bool   `json:"array"`     //是否数组
	Struct    Struct `json:"struct"`    //是结构体时
}

const (
	TagParamPath = "uri"
	TagParamFrom = "form"
	TagParamJson = "json"
	TagParamXml  = "xml"
	TagParamYaml = "yaml"

	OpenApiInPath   = "path"
	OpenApiInQuery  = "query"
	OpenApiInHeader = "header"
	OpenApiInCookie = "cookie"

	TagValidate = "validate"

	OpenApiTypeArray   = "array"
	OpenApiTypeBoolean = "boolean"
	OpenApiTypeInteger = "integer"
	OpenApiTypeNull    = "null"
	OpenApiTypeNumber  = "number"
	OpenApiTypeObject  = "object"
	OpenApiTypeString  = "string"

	OpenApiSchemasPrefix = "#/components/schemas/"
)

var tagParams = []string{TagParamPath, TagParamFrom, TagParamJson, TagParamXml, TagParamYaml}

func GetTagInfo(fieldTag string) (tag, name, validate string) {
	if fieldTag == "" {
		return
	}
	st := reflect.StructTag(fieldTag[1 : len(fieldTag)-1])
	for _, tagParam := range tagParams {
		val := st.Get(tagParam)
		if val != "" {
			tag = tagParam
			if tag == TagParamPath {
				tag = OpenApiInPath
			}
			name = util.GetJsonNameFromTag(val)
		}
	}
	validate = st.Get(TagValidate)
	return
}

var integerTypes = BaseType{"uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "int", "byte", "rune"}
var boolTypes = BaseType{"bool"}
var numberTypes = BaseType{"float32", "float64"}
var stringTypes = BaseType{"string"}

var baseTypes = BaseType{"interface{}", "error", "any", "bool", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "string", "int", "byte", "rune"}

type BaseType []string

func (bt BaseType) CheckIn(typ string) bool {
	for _, baseType := range bt {
		if baseType == typ || baseType == "[]"+baseType {
			return true
		}
	}
	return false
}

func (f Field) ToProperty() []Property {
	var p []Property
	if !baseTypes.CheckIn(f.Type) {
		if f.Name == "" { //组合参数
			for _, field := range f.Struct.Fields {
				if field.Type != f.Struct.Name {
					p = append(p, field.ToProperty()...)
				} else {
					p = append(p, Property{Name: f.ParamName, Description: f.Comment, Type: PropertyTypeObject, Properties: make(map[string]Property)})
				}
			}
		} else { //对象
			property := Property{Name: f.ParamName, Description: f.Comment, Type: PropertyTypeObject, Properties: make(map[string]Property)}
			for _, field := range f.Struct.Fields {
				if field.Type != f.Struct.Name {
					for _, fp := range field.ToProperty() {
						property.Properties[fp.Name] = fp
					}
				} else {
					property.Properties[f.ParamName] = Property{Name: f.ParamName, Description: f.Comment, Type: PropertyTypeObject, Properties: make(map[string]Property)}
				}

			}
			p = append(p, property)
		}
	} else {
		p = append(p, Property{Name: f.ParamName, Description: f.Comment, Type: f.GetOpenApiType(), Format: f.Type})
	}
	return p
}

func (f Field) ToParameter() Parameter {
	param := Parameter{
		Name:        f.ParamName,
		In:          f.GetOpenApiIn(),
		Description: f.Comment,
		Schema:      Property{Type: f.GetOpenApiType(), Format: f.Type},
	}
	if f.Validate != "" {
		param.Required = true
		if f.Validate != "required" {
			param.Description += "参数验证:" + f.Validate
		}
	}
	return param
}

func (f *Field) GetOpenApiIn() string {
	if f.In == TagParamFrom {
		return OpenApiInQuery
	}
	return f.In
}

func (f *Field) GetOpenApiType() string {
	if integerTypes.CheckIn(f.Type) {
		return OpenApiTypeInteger
	} else if numberTypes.CheckIn(f.Type) {
		return OpenApiTypeNumber
	} else if boolTypes.CheckIn(f.Type) {
		return OpenApiTypeBoolean
	} else if stringTypes.CheckIn(f.Type) {
		return OpenApiTypeString
	} else {
		if strings.Index(f.Type, "[]") == 0 {
			return OpenApiTypeArray
		} else if !baseTypes.CheckIn(f.Type) {
			return OpenApiTypeObject
		}
	}
	return f.Type
}
