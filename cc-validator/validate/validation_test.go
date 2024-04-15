package validate

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidation(t *testing.T) {
	var (
		requiredOk     = func(input string) error { return nil }
		requiredFailed = func(input string) error { return fmt.Errorf("%s", input) }
	)

	t.Run("Required checks should fail on errors and Forbid checks should pass", func(t *testing.T) {
		gotValidation := &Validation{}
		gotValidation.Require(requiredOk)
		ok, errs := gotValidation.Validate("")
		assert.Lenf(t, errs, 0, "expecting 0 of 1 validation check to fail, got: %s", errors.Join(errs...))
		assert.True(t, ok, "validation check should have succeeded when Check function returns nil")

		gotValidation = &Validation{}
		gotValidation.Require(requiredFailed)
		ok, errs = gotValidation.Validate("validation-error-123")
		assert.Lenf(t, errs, 1, "expecting 1 error from validation check, got: %s", errors.Join(errs...))
		assert.ErrorContains(t, errors.Join(errs...), "validation-error-123", "the error from within the wrapped function should be accessible")
		assert.False(t, ok, "validation check should have failed when Check function returned an error")

	})
}
