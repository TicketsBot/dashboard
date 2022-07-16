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

func Slice[T any](v ...T) []T {
	return v
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
