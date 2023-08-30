package convert

func ToInt(i interface{}) int {
	switch dst := i.(type) {
	case int:
		return dst
	case int32:
		return int(dst)
	case int64:
		return int(dst)
	}
	panic("convert.ToInt: invalid type")
}

// ToInt64 转换为int64
func ToInt64[T ~int | ~int32 | ~int64](i T) int64 {
	return int64(i)
}

// ToInt32 转换为int32
func ToInt32[T ~int | ~int32 | ~int64](i T) int32 {
	return int32(i)
}
