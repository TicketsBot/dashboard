package defaults

type ShouldApplyCheck func() bool
type ApplicatorFunc func()

type DefaultApplicator struct {
	ShouldApply ShouldApplyCheck
	Apply       ApplicatorFunc
}

func NewDefaultApplicator[T any](shouldApplyGenerator func(T) ShouldApplyCheck, ptr *T, defaultValue T) DefaultApplicator {
	return DefaultApplicator{
		ShouldApply: shouldApplyGenerator(*ptr),
		Apply: func() {
			*ptr = defaultValue
		},
	}
}

func ApplyDefaults(applicators ...DefaultApplicator) {
	for _, applicator := range applicators {
		if applicator.ShouldApply() {
			applicator.Apply()
		}
	}
}
