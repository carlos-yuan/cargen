package set

type Empty struct{}

func MapToSplice[T comparable, V any](m map[T]V) []T {
	list := make([]T, len(m))
	i := 0
	for k := range m {
		list[i] = k
		i++
	}
	return list
}
