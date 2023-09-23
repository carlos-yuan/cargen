package convert

import (
	"reflect"
	"runtime"
	"strings"
)

func GetStructModAndName(s any) (mod, name string) {
	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	mod = typ.PkgPath()[:strings.Index(typ.PkgPath(), "/")]
	return mod, typ.Name()
}

func GetMethodName(method interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
	name = name[strings.LastIndex(name, ".")+1 : strings.LastIndex(name, "-")]
	return name
}
