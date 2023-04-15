package utils

func Ptr[T any](v T) *T {
	return &v
}

func ValueOrZero[T any](v *T) T {
	if v == nil {
		return *new(T)
	} else {
		return *v
	}
}

func SetNilIfZero[T comparable](value **T) {
	if value != nil && *value != nil && **value == *new(T) {
		*value = nil
	}
}

func Slice[T any](v ...T) []T {
	return v
}

func Exists[T comparable](v []T, el T) bool {
	for _, e := range v {
		if e == el {
			return true
		}
	}

	return false
}

func ExistsMap[T any, U comparable](v []T, el U, mapper func(T) U) bool {
	for _, e := range v {
		mapped := mapper(e)
		if mapped == el {
			return true
		}
	}

	return false
}

func FindMap[T any, U comparable](v []T, el U, mapper func(T) U) *T {
	for _, e := range v {
		mapped := mapper(e)
		if mapped == el {
			return &e
		}
	}

	return nil
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
