package validation

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCtx struct {
	a int
	b string
}

func validateGreaterThanZero(ctx testCtx) ValidationFunc {
	return func() error {
		if ctx.a <= 0 {
			return NewInvalidInputError("a must be greater than 0")
		}

		return nil
	}
}

func validateStringNotEmpty(ctx testCtx) ValidationFunc {
	return func() error {
		if len(ctx.b) == 0 {
			return NewInvalidInputError("b must be greater than 0")
		}

		return nil
	}
}

func TestSuccessful(t *testing.T) {
	ctx := testCtx{a: 1, b: "test"}

	if err := Validate(context.Background(), ctx, validateGreaterThanZero, validateStringNotEmpty); err != nil {
		t.Error(err)
	}
}

func TestNoValidators(t *testing.T) {
	ctx := testCtx{a: 1, b: "test"}

	if err := Validate(context.Background(), ctx); err != nil {
		t.Error(err)
	}
}

func TestSingleFail(t *testing.T) {
	ctx := testCtx{a: 1, b: ""}
	err := Validate(context.Background(), ctx, validateGreaterThanZero, validateStringNotEmpty)
	if err == nil {
		t.Fatal("expected error")
	}

	var validationError *InvalidInputError
	if !errors.As(err, &validationError) {
		t.Fatal("expected InvalidInputError error")
	}

	assert.Equal(t, "b must be greater than 0", validationError.Message)
}

func TestDualFail(t *testing.T) {
	ctx := testCtx{a: 0, b: ""}
	err := Validate(context.Background(), ctx, validateGreaterThanZero, validateStringNotEmpty)
	if err == nil {
		t.Error("expected error")
	}

	var validationError *InvalidInputError
	if !errors.As(err, &validationError) {
		t.Error("expected InvalidInputError error")
	}

	if validationError.Message != "a must be greater than 0" && validationError.Message != "b must be greater than 0" {
		t.Errorf("got wrong error message: %s", validationError.Message)
	}
}
