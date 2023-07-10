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

func FistIsLower(str string) bool {
	if len(str) == 0 {
		return false
	}
	word := str[0]
	return 'a' <= word && word <= 'z'
}
