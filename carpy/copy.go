package carpy

import (
	"bytes"
	"fmt"
	openapi "github.com/carlos-yuan/cargen/open_api"
	"github.com/carlos-yuan/cargen/util/md5"
	"go/ast"
	"strings"
)

func Gen(base string) {
	cp := Carpy{}
	cp.Gen(base)
}

const (
	PackageName   = "carpy"
	InterfaceName = "Copy"
)

type Carpy struct {
	cpPkg        map[string]openapi.Package
	pkgs         *openapi.Packages
	cpPkgStructs map[string]*copyStructInfoList
}

// Gen 生成
func (c *Carpy) Gen(base string) {
	c.pkgs = &openapi.Packages{}
	c.pkgs.InitPackages(base)
	c.findCopyPkg()
	c.findCopyStruct()
	c.generateCopyFile()
}

// findCopyPkg 查找包下所有拷贝信息
func (c *Carpy) findCopyPkg() {
	c.cpPkg = make(map[string]openapi.Package)
	for _, pkg := range *c.pkgs {
		for _, file := range pkg.GetAstPkg().Files {
			if file.Scope == nil {
				continue
			}
			for name, obj := range file.Scope.Objects {
				if d, ok := obj.Decl.(*ast.ValueSpec); ok {
					field := openapi.GetExprInfo(d.Type)
					if field.Pkg == PackageName && field.Type == InterfaceName {
						c.cpPkg[name] = pkg
					}
				}
			}
		}
	}
}

// findCopyStruct 查找所有需要复制的类型
func (c *Carpy) findCopyStruct() {
	c.cpPkgStructs = make(map[string]*copyStructInfoList)
	for name, pkg := range c.cpPkg {
		for _, file := range pkg.GetAstPkg().Files {
			vsr := c.newVisitor(name, openapi.FindImports(file))
			ast.Walk(vsr, file)
			for _, info := range vsr.structs {
				for _, from := range info.from {
					if c.cpPkgStructs[name] == nil {
						c.cpPkgStructs[name] = &copyStructInfoList{}
					}
					c.cpPkgStructs[name].append(info.to, from)
				}
			}
		}
	}
}

func (c *Carpy) generateCopyFile() {
	for name, structs := range c.cpPkgStructs {
		pkg := c.cpPkg[name]
		imports := make(map[string]string) // map[包路径]包名+md5(包路径) 用于区别包名一样路径不一样的包
		for _, st := range *structs {
			if pkg.Path != st.to.Pkg.Path { //当前包不处理名称
				imports[st.to.Pkg.Path] = st.to.Pkg.Name + md5.Encode16(st.to.Pkg.Path)
			}
			for _, from := range st.from {
				if pkg.Path != from.Pkg.Path { //当前包不处理名称
					imports[from.Pkg.Path] = from.Pkg.Name + md5.Encode16(from.Pkg.Path)
				}
			}
		}
		var buf bytes.Buffer
		var importBuf bytes.Buffer
		for path, name := range imports {
			importBuf.WriteString(fmt.Sprintf("%s \"%s\"\n", name, path))
		}
		buf.WriteString(fmt.Sprintf(templateHeader, pkg.Name, importBuf.String()))
		copyStructName := "carpy" + pkg.Name
		var caseBuf bytes.Buffer
		var funcBuf bytes.Buffer
		for _, st := range *structs {
			to := st.to
			toName := to.Name
			if pkg.Name != to.Pkg.Name || pkg.Path != to.Pkg.Path {
				toName = imports[to.Pkg.Path] + "." + to.Name //包名.类型名
			}
			var caseFromBuf bytes.Buffer
			for funName, from := range st.from {
				fromName := from.Name
				if pkg.Name != from.Pkg.Name || pkg.Path != from.Pkg.Path {
					fromName = imports[from.Pkg.Path] + "." + from.Name //包名.类型名
				}
				var funcBody bytes.Buffer
				caseFromBuf.WriteString(fmt.Sprintf(templateFromCase, fromName, funName))
				for _, tf := range to.Fields {
					for _, ff := range from.Fields {
						if tf.Name == ff.Name && (tf.Name[0] > 64 && tf.Name[0] < 91 || (!strings.Contains(toName, ".") && !strings.Contains(fromName, "."))) { //可导出或本包不可导出才能拷贝
							if tf.Type == ff.Type {
								funcBody.WriteString(fmt.Sprintf("\tto.%s = from.%s\n", tf.Name, ff.Name))
							} else if (tf.Type == "int" || tf.Type == "int64") && ( //以下为整形及浮点的类型强转
							ff.Type == "int" || ff.Type == "int64" || ff.Type == "int32" || ff.Type == "rune" || ff.Type == "int16" || ff.Type == "int8" ||
								ff.Type == "uint32" || ff.Type == "uint16" || ff.Type == "uint8") {
								funcBody.WriteString(fmt.Sprintf("\tto.%s = %s(from.%s)\n", tf.Name, tf.Type, ff.Name))
							} else if (tf.Type == "int32") && (ff.Type == "rune" || ff.Type == "int16" || ff.Type == "int8" ||
								ff.Type == "uint16" || ff.Type == "uint8") {
								funcBody.WriteString(fmt.Sprintf("\tto.%s = %s(from.%s)\n", tf.Name, tf.Type, ff.Name))
							} else if tf.Type == "int16" && (ff.Type == "int8" || ff.Type == "uint8") {
								funcBody.WriteString(fmt.Sprintf("\tto.%s = %s(from.%s)\n", tf.Name, tf.Type, ff.Name))
							} else if (tf.Type == "uint" || tf.Type == "uint64") && (ff.Type == "uint" || ff.Type == "uint64" ||
								ff.Type == "uint32" || ff.Type == "uint16" || ff.Type == "uint8") {
								funcBody.WriteString(fmt.Sprintf("\tto.%s = %s(from.%s)\n", tf.Name, tf.Type, ff.Name))
							} else if tf.Type == "uint32" && (ff.Type == "uint16" || ff.Type == "uint8") {
								funcBody.WriteString(fmt.Sprintf("\tto.%s = %s(from.%s)\n", tf.Name, tf.Type, ff.Name))
							} else if tf.Type == "uint16" && ff.Type == "uint8" {
								funcBody.WriteString(fmt.Sprintf("\tto.%s = %s(from.%s)\n", tf.Name, tf.Type, ff.Name))
							} else if (tf.Type == "float64" || tf.Type == "float32") && (ff.Type == "float32" || ff.Type == "int" || ff.Type == "int64" || ff.Type == "int32" || ff.Type == "rune" || ff.Type == "int16" || ff.Type == "int8" ||
								ff.Type == "uint64" || ff.Type == "uint32" || ff.Type == "uint16" || ff.Type == "uint8") {
								funcBody.WriteString(fmt.Sprintf("\tto.%s = %s(from.%s)\n", tf.Name, tf.Type, ff.Name))
							}
						}
					}
				}
				funcBuf.WriteString(fmt.Sprintf(templateFun, funName, toName, fromName, funcBody.String()))
			}
			caseBuf.WriteString(fmt.Sprintf(templateToCase, toName, caseFromBuf.String()))
		}
		buf.WriteString(fmt.Sprintf(templateDecl, name, copyStructName, copyStructName, copyStructName, caseBuf.String()))
		buf.WriteString(funcBuf.String())
		println(buf.String())
		println("123")
	}
}

const templateHeader = `
package %s

import (
	"errors"
	"github.com/carlos-yuan/cargen/carpy"
%s
)

`

const templateDecl = `func init() {
	%s = &%s{}
}

type %s struct{}

func (c *%s) Copy(to any, from any, opts ...carpy.CopyOption) error {
	if to == nil || from == nil {
		return nil
	}
	switch to := to.(type) {
	%s
	default:
		return errors.New("unknown copy to " + carpy.GetTypeName(to))
	}
}
`

const templateToCase = `
	case *%s:
		switch from := from.(type) {
%s
		default:
			return errors.New("unknown copy from " + carpy.GetTypeName(from))
		}

`
const templateFromCase = `
		case *%s:
			return %s(to, from, opts...)
`

const templateFun = `
func %s(to *%s, from *%s, opts ...carpy.CopyOption) (err error) {
%s
	return
}

`

const templateOption = `
	for _, opt := range opts {
		dst, err := opt(to.Count, from.Count)
		if err != nil {
			return err
		}
		if count, ok := dst.(int32); ok {
			to.Count = count
		}
	}
`

type copyStructInfo struct {
	to   *openapi.Struct
	from map[string]*openapi.Struct
}

type visitor struct {
	name    string
	carpy   *Carpy
	structs copyStructInfoList
	imports map[string]string
}

type copyStructInfoList []*copyStructInfo

func (c *copyStructInfoList) append(to *openapi.Struct, from *openapi.Struct) {
	funcName := "Copy" + from.Name + "To" + to.Name + md5.Encode16(to.Pkg.Path+from.Pkg.Path)
	if len(*c) == 0 {
		*c = append(*c, &copyStructInfo{to: to, from: map[string]*openapi.Struct{funcName: from}})
	} else {
		find := false
		for i := range *c {
			if (*c)[i].to.Name == to.Name && (*c)[i].to.Pkg.Name == to.Pkg.Name && (*c)[i].to.Pkg.Path == to.Pkg.Path {
				(*c)[i].from[funcName] = from
				find = true
				break
			}
		}
		if !find {
			*c = append(*c, &copyStructInfo{to: to, from: map[string]*openapi.Struct{funcName: from}})
		}
	}
}

func (c *Carpy) newVisitor(name string, imports map[string]string) *visitor {
	return &visitor{name: name, carpy: c, imports: imports, structs: copyStructInfoList{}}
}

// 寻找所需的
func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if call, ok := n.(*ast.CallExpr); ok {
		if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selector.X.(*ast.Ident); ok {
				if ident.Name == v.name && selector.Sel.Name == InterfaceName {
					to, from := v.GetCallArgsInfo(call.Args)
					if to != nil && from != nil {
						v.structs.append(to, from)
					}
				}
			}
		}
	}
	if n == nil {
		return nil
	}
	return v
}

func (v *visitor) GetCallArgsInfo(args []ast.Expr) (to *openapi.Struct, from *openapi.Struct) {
	if len(args) < 2 {
		return
	}
	for i, arg := range args {
		if i > 1 {
			break
		}
		info := getExprInfo(arg)
		pkg := v.carpy.cpPkg[v.name]
		path := pkg.Path
		if info.Pkg != "" {
			path = v.imports[info.Pkg]
		} else {
			info.Pkg = pkg.Name
		}
		if i == 0 {
			to = v.carpy.pkgs.FindStructPtr(path, info.Pkg, info.Type)
			if to == nil {
				return
			}
		} else if i == 1 {
			from = v.carpy.pkgs.FindStructPtr(path, info.Pkg, info.Type)
			if from == nil {
				return
			}
		}
	}
	return
}

func getExprInfo(expr ast.Expr) *openapi.Field {
	switch exp := expr.(type) {
	case *ast.Ident:
		field := &openapi.Field{Name: exp.Name}
		f := GetObject(exp.Obj)
		if f != nil {
			field.Pkg = f.Pkg
			field.Type = f.Type
		}
		return field
	case *ast.SelectorExpr:
		return &openapi.Field{Type: exp.Sel.Name, Pkg: exp.X.(*ast.Ident).Name}
	case *ast.StarExpr:
		return getExprInfo(exp.X)
	case *ast.ArrayType:
		f := getExprInfo(exp.Elt)
		f.Array = true
		return f
	case *ast.SliceExpr:
		return getExprInfo(exp.X)
	case *ast.UnaryExpr:
		return getExprInfo(exp.X)
	case *ast.CompositeLit:
		return GetCompositeLitInfo(exp)
	}
	return nil
}

func GetObject(obj *ast.Object) *openapi.Field {
	if obj == nil || obj.Decl == nil {
		return nil
	}
	switch decl := obj.Decl.(type) {
	case *ast.ValueSpec:
		return getExprInfo(decl.Values[0])
	}
	return nil
}

func GetCompositeLitInfo(obj *ast.CompositeLit) *openapi.Field {
	if obj == nil || obj.Type == nil {
		return nil
	}
	f := getExprInfo(obj.Type)
	//需要转换为类型
	if _, ok := obj.Type.(*ast.Ident); ok && f != nil {
		f.Type = f.Name
	}
	return f
}
