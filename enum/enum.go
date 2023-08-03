package enum

import (
	"bytes"
	"fmt"
	"github.com/carlos-yuan/cargen/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
)

func GenEnum(path, dictTable, dictType, dictName, dictLabel, dictValue, dsn string) {
	db, _ := gorm.Open(mysql.Open(dsn))
	var dict []map[string]interface{}
	err := db.Table(dictTable).Find(&dict).Error
	if err != nil {
		println(err)
	}
	sort.Slice(dict, func(i, j int) bool {
		if dict[i] == nil && dict[j] == nil {
			istr := dict[i][dictType].(string) + dict[i][dictName].(string)
			jstr := dict[j][dictType].(string) + dict[j][dictName].(string)
			return strings.Compare(istr, jstr) == -1
		}
		return false
	})
	var buf bytes.Buffer
	buf.WriteString("package enum\n")
	dicts := make(map[string][]Dict)
	for _, m := range dict {
		typ := util.ToCamelCase(strings.TrimSpace(m[dictType].(string)))
		dicts[typ] = append(dicts[typ], Dict{Type: typ, Name: util.ToCamelCase(strings.TrimSpace(m[dictName].(string))), Label: strings.TrimSpace(m[dictLabel].(string)), Value: strings.TrimSpace(m[dictValue].(string))})
	}
	keys := util.MapToSplice(dicts)
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) == -1
	})
	for _, key := range keys {
		list := dicts[key]
		//之前类型生成代码
		var constCode, caseCode string
		for _, d := range list {
			if d.Name == "" {
				continue
			}
			val, err := strconv.ParseInt(d.Value, 10, 64)
			if err != nil {
				panic(err)
			}
			constName := d.Type + d.Name
			constCode += fmt.Sprintf(enumConstTemplate, constName, val)
			caseCode += fmt.Sprintf(enumCaseTemplate, constName, strings.ReplaceAll(d.Label, `"`, `\"`))
		}
		if len(constCode) > 0 {
			buf.WriteString(fmt.Sprintf(enumTemplate, key, constCode, key, caseCode))
		}
	}
	err = util.WriteByteFile(path+"/enum/enum.go", []byte(buf.String()))
	if err != nil {
		panic(err)
	}
}

type Dict struct {
	Name  string
	Type  string
	Label string
	Value string
}

const enumConstTemplate = `	%s = %d
`

const enumCaseTemplate = `
	case %s:
		return "%s"`

const enumTemplate = `
type %s int

const (
%s
)

func (t %s) String() string {
	switch t {
%s
	}
	return ""
}
`
