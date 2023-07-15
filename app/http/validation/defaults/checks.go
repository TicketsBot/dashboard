package defaults

func EmptyStringCheck(s string) ShouldApplyCheck {
	return func() bool {
		return len(s) == 0
	}
}

func NilCheck[T any](v *T) ShouldApplyCheck {
	return func() bool {
		return v == nil
	}
}

func NilOrEmptyStringCheck(s *string) ShouldApplyCheck {
	return func() bool {
		return s == nil || len(*s) == 0
	}
}
