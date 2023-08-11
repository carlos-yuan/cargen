package convert

import "strings"

func Hide(str string, length int) string {
	max := len(str)
	bytes := []uint8(str)
	if max-1 < length {
		builder := strings.Builder{}
		builder.WriteByte(bytes[0])
		for i := 0; i < length; i++ {
			builder.WriteRune(42)
		}
		if max > 1 {
			builder.WriteByte(bytes[max-1])
		}
		return builder.String()
	} else {
		start := (max - length) / 2
		for i := 0; i < max; i++ {
			if i >= start && i < start+length {
				bytes[i] = uint8(42)
			}
		}
		return string(bytes)
	}
}

func HideCompare(hideStr, source string) bool {
	fir, sec := "", ""
	start := strings.Index(hideStr, "*")
	if start != -1 {
		fir = hideStr[:start]
	}
	end := strings.LastIndex(hideStr, "*")
	if end != -1 {
		sec = hideStr[end+1:]
	}
	if fir != "" {
		if strings.Index(source, fir) != 0 {
			return false
		}
	}
	if sec != "" {
		if strings.LastIndex(source, sec) != len(source)-len(sec) {
			return false
		}
	}
	return true
}
