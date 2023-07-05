package util

func FistToLower(str string) string {
	if len(str) == 0 {
		return str
	}
	word := str[0]
	if 'A' <= word && word <= 'Z' {
		word += 'a' - 'A'
	}
	byt := []byte(str)
	byt[0] = word
	return string(byt)
}
