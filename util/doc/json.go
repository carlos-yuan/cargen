package doc

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

func ToJsonString(js interface{}) string {
	b, _ := json.Marshal(js)
	if b != nil {
		return string(b)
	}
	return ""
}

func PrintJson(js interface{}) {
	b, _ := json.Marshal(js)
	if b != nil {
		println("json", string(b))
	}
}

func MapToAssciiSortJson(m map[string]interface{}) string {
	i := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i, k := range keys {
		var val string
		switch m[k].(type) {
		case map[string]interface{}:
			val = MapToAssciiSortJson(m[k].(map[string]interface{}))
		case []interface{}:
			var arrBuffer bytes.Buffer
			arrBuffer.WriteString("[")
			for idx, sv := range m[k].([]interface{}) {
				switch sv.(type) {
				case map[string]interface{}:
					arrBuffer.WriteString(MapToAssciiSortJson(m[k].([]interface{})[idx].(map[string]interface{})))
				case string:
					arrBuffer.WriteString(`"`)
					arrBuffer.WriteString(m[k].([]interface{})[idx].(string))
					arrBuffer.WriteString(`"`)
				case int:
					arrBuffer.WriteString(strconv.Itoa(m[k].([]interface{})[idx].(int)))
				case int32:
					arrBuffer.WriteString(strconv.Itoa((int)(m[k].([]interface{})[idx].(int32))))
				case int64:
					arrBuffer.WriteString(strconv.Itoa((int)(m[k].([]interface{})[idx].(int64))))
				}
				if idx != len(m[k].([]interface{}))-1 {
					arrBuffer.WriteString(",")
				}
			}
			arrBuffer.WriteString("]")
			val = arrBuffer.String()
		case []string, []int, []int32, []int64:
			b, _ := json.Marshal(m[k])
			val = string(b)
		case string:
			val = `"` + m[k].(string) + `"`
		case int:
			val = strconv.Itoa(m[k].(int))
		case int32:
			val = strconv.Itoa((int)(m[k].(int32)))
		case int64:
			val = strconv.Itoa((int)(m[k].(int64)))
		}
		buffer.WriteString(`"`)
		buffer.WriteString(k)
		buffer.WriteString(`":`)
		buffer.WriteString(val)
		if i != len(keys)-1 {
			buffer.WriteString(`,`)
		}
	}
	buffer.WriteString("}")
	return buffer.String()
}

func GetTagName(tag string) string {
	if strings.Index(tag, ",") != -1 {
		return tag[:strings.Index(tag, ",")]
	}
	return tag
}
