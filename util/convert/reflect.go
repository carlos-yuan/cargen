package convert

import (
	"reflect"
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
