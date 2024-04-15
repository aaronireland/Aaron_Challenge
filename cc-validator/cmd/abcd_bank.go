package cmd

import (
	v "github.com/aaronireland/Aaron_Challenge/cc-validator/validate"
)

type abcdBankValidator interface {
	Validate(input string) (validated bool, errs []error)
}

func NewAbcdBankValidator() abcdBankValidator {
	validator := &v.Validation{}
	validator.Require(v.LeadingDigit(4, 5, 6))
	validator.Require(v.DigitCountEquals(16))
	validator.Require(v.OnlyDigitsOrValidSeparator("-", 4))
	validator.Require(v.DoesNotExceedMaxRepeatingDigits(3))

	return validator
}
