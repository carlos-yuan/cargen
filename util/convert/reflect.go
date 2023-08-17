package convert

import (
	"reflect"
	"strings"
)

func GetStructModAndName(s any) (mod, name string) {
	typ := reflect.TypeOf(s).Elem()
	mod = typ.PkgPath()[:strings.Index(typ.PkgPath(), "/")]
	return mod, typ.Name()
}
