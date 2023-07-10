package openapi

type Struct struct {
	Name      string            `json:"name"`      //名称
	Des       string            `json:"des"`       //描述
	Type      string            `json:"type"`      //结构体类型
	Field     string            `json:"field"`     //所需字段 有时可能需要的是结构体中的字段 如：rsp.List
	Fields    []Field           `json:"fields"`    //字段
	MethodMap *MethodMap        `json:"methodMap"` //map[参数位置][]链路 方法的返回值类型查找
	Imports   map[string]string //包下所有导入的包信息
	Methods   []StructMethod    `json:"methods"` //方法列表 包括接口的方法
	Api       []Api             //所有API接口信息
	pkg       *Package
}

type MethodMap struct {
	Idx   int
	Paths []string
}

func NewStruct(pkg *Package) Struct {
	return Struct{pkg: pkg}
}

type StructMethod struct {
	Name    string  //名称
	Pkg     string  //包名
	PkgPath string  //包路径
	Args    []Field //参数
	Returns []Field //返回值
}

func (sct *Struct) GetField() *Field {
	for _, field := range sct.Fields {
		if field.Name == sct.Field {
			sct.Type = field.Type
			return &field
		}
	}
	return nil
}
