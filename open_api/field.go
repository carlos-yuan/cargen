package openapi

import (
	"github.com/carlos-yuan/cargen/util"
	"reflect"
)

type Field struct {
	Name      string  `json:"name"`      //名称
	Type      string  `json:"type"`      //类型
	Tag       string  `json:"tag"`       //tag
	In        string  `json:"in"`        //存在类型 path 路径 json json对象内
	ParamName string  `json:"paramName"` //参数名
	Validate  string  `json:"validate"`  //验证
	Pkg       string  `json:"pkg"`       //包名 类型为结构体时
	PkgPath   string  `json:"pkgPath"`   //包路径 类型为结构体时
	Comment   string  `json:"comment"`   //注释
	Array     bool    `json:"array"`     //是否数组
	Struct    *Struct `json:"struct"`    //是结构体时
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

var baseTypes = BaseType{"error", "bool", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "string", "int", "byte", "rune"}

type BaseType []string

func (bt BaseType) CheckIn(typ string) bool {
	for _, baseType := range bt {
		if baseType == typ || baseType == "[]"+baseType {
			return true
		}
	}
	return false
}

// ToProperty 找到属性
// deep 递归层数 max递归深度
func (f Field) ToProperty(storey, deep int) []Property {
	if f.Name != "" && util.FistIsLower(f.Name) { //小写开头的隐藏字段去掉
		return []Property{}
	}
	if storey > deep { //深度限制
		return []Property{}
	}
	var p []Property
	if !baseTypes.CheckIn(f.Type) { //非一般类型
		if f.Name == "" { //组合参数T{K}
			for _, field := range f.Struct.Fields {
				p = append(p, field.ToProperty(storey+1, deep)...)
			}
		} else { //嵌套对象
			property := Property{Name: f.ParamName, Description: f.Comment, Properties: make(map[string]Property)}
			if f.Array {
				property.Type = PropertyTypeArray
			} else {
				property.Type = PropertyTypeObject
			}
			if f.Struct != nil {
				var fpps []Property
				for _, field := range f.Struct.Fields {
					pps := field.ToProperty(storey+1, deep)
					if len(pps) == 0 {
						continue
					}
					p := Property{Name: field.ParamName, Description: field.Comment, Properties: make(map[string]Property)}
					if len(pps) == 1 {
						p = pps[0]
					} else {
						if !field.Array {
							if len(pps) == 1 && baseTypes.CheckIn(pps[0].Type) {
								p = pps[0]
							} else {
								for _, fp := range pps {
									p.Properties[fp.Name] = fp
								}
								p.Type = PropertyTypeObject
							}
						} else {
							pp := Property{Properties: make(map[string]Property)}
							if len(pps) == 1 && baseTypes.CheckIn(pps[0].Type) {
								pp = pps[0]
							} else {
								for _, fp := range pps {
									pp.Properties[fp.Name] = fp
								}
								pp.Type = PropertyTypeObject
							}
							p.Items = &pp
							p.Type = PropertyTypeArray
						}
					}
					fpps = append(fpps, p)
				}
				if len(fpps) == 1 {

				}
				if !f.Array {
					if len(fpps) == 1 && baseTypes.CheckIn(fpps[0].Type) {
						p = append(p, fpps[0])
					} else {
						for _, fp := range fpps {
							property.Properties[fp.Name] = fp
						}
						property.Type = PropertyTypeObject
						p = append(p, property)
					}
				} else {
					pp := Property{Properties: make(map[string]Property)}
					if len(fpps) == 1 && baseTypes.CheckIn(fpps[0].Type) {
						pp = fpps[0]
					} else {
						for _, fp := range fpps {
							pp.Properties[fp.Name] = fp
						}
					}
					property.Items = &pp
					property.Type = PropertyTypeArray
					p = append(p, property)
				}
			}
		}
	} else { //一般类型
		if f.Array { //数组类型
			pp := Property{Name: f.ParamName, Description: f.Comment, Type: PropertyTypeArray, Format: f.Type}
			pp.Items = &Property{Type: f.GetOpenApiType()}
			p = append(p, pp)
		} else {
			p = append(p, Property{Name: f.ParamName, Description: f.Comment, Type: f.GetOpenApiType(), Format: f.Type})
		}
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
		if f.Array {
			return OpenApiTypeArray
		} else if !baseTypes.CheckIn(f.Type) {
			return OpenApiTypeObject
		}
	}
	return f.Type
}

type Fields []Field

// ToProperty 找到属性
// deep 递归层数 max递归深度
func (f *Fields) ToProperty(storey, deep int) []Property {
	var list []Property
	for _, field := range *f {
		pps := field.ToProperty(storey, deep)
		if baseTypes.CheckIn(field.Type) && len(pps) == 1 {
			list = append(list, pps[0])
		} else {
			p := Property{Name: field.ParamName, Description: field.Comment, Properties: make(map[string]Property)}
			if !field.Array {
				for _, fp := range pps {
					p.Properties[fp.Name] = fp
				}
				p.Type = PropertyTypeObject
			} else {
				pp := Property{Properties: make(map[string]Property)}
				if len(pps) == 1 && baseTypes.CheckIn(pps[0].Type) {
					pp = pps[0]
				} else {
					for _, fp := range pps {
						pp.Properties[fp.Name] = fp
					}
					pp.Type = PropertyTypeObject
				}
				p.Items = &pp
				p.Type = PropertyTypeArray
			}
			list = append(list, p)
		}
	}
	return list
}
