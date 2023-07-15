package validation

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type Validator[T any] func(validationContext T) ValidationFunc
type ValidationFunc func() error

func Validate[T any](ctx context.Context, validationContext T, validators ...Validator[T]) error {
	group, _ := errgroup.WithContext(ctx)

	for _, validator := range validators {
		validator := validator
		group.Go(validator(validationContext))
	}

	return group.Wait()
}
