package carpy

import "reflect"

// Copy 用于拷贝的接口 在生成文件中将生成该接口的实现体
type Copy interface {
	Copy(to any, from any, opts ...CopyOption) error
}

func GetTypeName(obj any) string {
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.Name()
}

// CopyOption 拷贝类型转换
type CopyOption func(to any, from any) (dst any, err error)
