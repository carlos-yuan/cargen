package openapi

import (
	"go/ast"

	"github.com/carlos-yuan/cargen/util/convert"
)

type Struct struct {
	Name      string            `json:"name"`      //名称
	Des       string            `json:"des"`       //描述
	Type      string            `json:"type"`      //结构体类型
	Field     string            `json:"field"`     //所需字段 有时可能需要的是结构体中的字段 如：rsp.List
	Fields    Fields            `json:"fields"`    //字段
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

func (sct Struct) Copy() Struct {
	s := Struct{
		Name:      sct.Name,
		Des:       sct.Name,
		Type:      sct.Type,
		Field:     sct.Field,
		Fields:    make([]Field, 0, len(sct.Fields)),
		MethodMap: sct.MethodMap,
		Imports:   sct.Imports,
		Methods:   sct.Methods,
		Api:       sct.Api,
		pkg:       sct.pkg,
	}
	for _, field := range sct.Fields {
		s.Fields = append(s.Fields, field)
	}
	return s
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

func (sct *Struct) GetStructFromAstInterfaceType(s *ast.InterfaceType) {
	for _, m := range s.Methods.List {
		sct.Methods = append(sct.Methods, sct.GetInterfaceMethodField(m))
	}
}

func (sct *Struct) FieldFromInterfaceAstField(fd *ast.Field) Field {
	f := GetFieldInfo(fd)
	if f.Pkg != "" {
		f.PkgPath = sct.Imports[f.Pkg]
	}
	if fd.Tag != nil {
		f.Tag = fd.Tag.Value
		f.In, f.ParamName, f.Validate = GetTagInfo(f.Tag)
	}
	f.Comment = FormatComment(fd.Comment)
	return f
}

// GetInterfaceMethodField 寻找接口字段定义
func (sct *Struct) GetInterfaceMethodField(it *ast.Field) StructMethod {
	sm := StructMethod{}
	if len(it.Names) == 1 {
		sm.Name = it.Names[0].Name
	}
	// todo 这点过滤掉了组合接口 基本上不会存在需要返回值时，是在组合接口中
	if ft, ok := it.Type.(*ast.FuncType); ok {
		if ft.Params != nil {
			for _, p := range ft.Params.List {
				sm.Args = append(sm.Args, sct.FieldFromAstField(p))
			}
		}
		if ft.Results != nil {
			for _, r := range ft.Results.List {
				sm.Returns = append(sm.Returns, sct.FieldFromAstField(r))
			}
		}
	}
	return sm
}

func (sct *Struct) GetStructFromAstStructType(s *ast.StructType) {
	if s != nil {
		for _, fd := range s.Fields.List {
			f := sct.FieldFromAstField(fd)
			if !baseTypes.CheckIn(f.Type) {
				field := GetExprInfo(fd.Type)
				if field.Type == ExprStruct {
					if f.Struct == nil {
						f.Struct = &Struct{}
					}
					f.Struct.Type = field.Pkg
					if st, ok := fd.Type.(*ast.StructType); ok {
						f.Struct.GetStructFromAstStructType(st)
					} else if arraytype, ok := fd.Type.(*ast.ArrayType); ok {
						if st, ok := arraytype.Elt.(*ast.StructType); ok {
							f.Struct.GetStructFromAstStructType(st)
						}
					}
				} else {
					f.Array = field.Array
					f.Type = field.Type
					f.Pkg = field.Pkg
					if f.Pkg != "" {
						f.PkgPath = sct.Imports[f.Pkg]
					} else {
						f.Pkg = sct.pkg.Name
						f.PkgPath = sct.pkg.Path
					}
					f.Struct = sct.pkg.pkgs.FindStructPtr(f.PkgPath, f.Pkg, f.Type)
				}
			}
			sct.Fields = append(sct.Fields, f)
		}
	}
}

// GetStructMethodFuncType 寻找方法参数定义
func (sct *Struct) GetStructMethodFuncType(ft *ast.FuncType) StructMethod {
	sm := StructMethod{}
	if ft.Params != nil {
		for _, p := range ft.Params.List {
			sm.Args = append(sm.Args, sct.FieldFromAstField(p))
		}
	}
	if ft.Results != nil {
		for _, r := range ft.Results.List {
			sm.Returns = append(sm.Returns, sct.FieldFromAstField(r))
		}
	}
	return sm
}

// FieldFromAstField 寻找基础的结构体字段定义
func (sct *Struct) FieldFromAstField(fd *ast.Field) Field {
	f := GetFieldInfo(fd)
	if f.Pkg != "" {
		f.PkgPath = sct.Imports[f.Pkg]
	}
	if fd.Tag != nil {
		f.Tag = fd.Tag.Value
		f.In, f.ParamName, f.Validate = GetTagInfo(f.Tag)
	}
	f.Comment = FormatComment(fd.Comment)
	return f
}

// ToProperty 找到属性
// deep 递归层数 max递归深度
func (sct Struct) ToProperty() Property {
	var pps []Property
	fieldArray := false
	for _, field := range sct.Fields {
		if sct.Field != "" {
			if field.Name == sct.Field {
				if field.Struct != nil {
					fieldArray = field.Array
					sct = *field.Struct
				}
			}
		}
	}
	for _, field := range sct.Fields {
		if field.Name != "" && convert.FistIsLower(field.Name) { //小写开头的隐藏字段去掉
			continue
		}
		pps = append(pps, field.ToProperty(0, 5)...)

	}
	var p = Property{Properties: make(map[string]Property)}
	if !fieldArray {
		for _, pp := range pps {
			p.Properties[pp.Name] = pp
		}
		p.Type = PropertyTypeObject
	} else {
		pp := Property{Type: PropertyTypeObject, Properties: make(map[string]Property)}
		for _, fpp := range pps {
			pp.Properties[fpp.Name] = fpp
		}
		p.Items = &pp
		p.Type = PropertyTypeArray
	}
	return p
}
