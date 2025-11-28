package collection

// Filter filters a collection.
func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// First finds the first item in the collection.
func First[T any](ss []T, test func(T) bool) (ret T) {
	for _, s := range ss {
		if test(s) {
			ret = s
		}
	}
	return
}

// Map a slice
func Map[T any, R any](ss []T, m func(T) R) (ret []R) {
	for _, s := range ss {
		ret = append(ret, m(s))
	}
	return
}

// Has checks if the collection has the item.
func Has[T any](ss []T, test func(T) bool) (ret bool) {
	for _, s := range ss {
		if test(s) {
			ret = true
		}
	}
	return
}

func Reduce[T any, R any](ss []T, fn func(R, T) R, initial R) (ret R) {
	ret = initial
	for _, v := range ss {
		ret = fn(ret, v)
	}
	return
}

func ReduceMap[K comparable, T any, R any](ss map[K]T, fn func(T, K, R) R, initial R) (ret R) {
	ret = initial
	for k, v := range ss {
		ret = fn(v, k, ret)
	}
	return
}
