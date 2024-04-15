package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeadingDigit(t *testing.T) {
	t.Run("Should return a validator function", func(t *testing.T) {
		check := LeadingDigit()

		if assert.NotNil(t, check) {

			// Valid leading digits based on default
			for _, test := range []string{
				"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "0O",
			} {
				err := check(test)
				assert.NoErrorf(t, err, "LeadingDigit check should accept any initial digit when no params are given, input was %s, got error: %s", test, err)
			}

			// Invalid leading digits
			for _, test := range []string{
				"a", "-", " 1", "_.4", "\\.1", "♳", "",
			} {
				err := check(test)
				assert.Errorf(t, err, "LeadingDigit failed to error on invalid leading digit: %s", test)
			}
		}

	})
}

func TestLeadingDigitWithParams(t *testing.T) {
	t.Run("Should return a validator function with expected regex", func(t *testing.T) {
		check := LeadingDigit(1, 3, 5, 7, 9)

		if assert.NotNil(t, check) {

			// Valid leading digits based on default
			for _, test := range []string{"1", "3", "5", "7", "9", "920", "7a52-"} {
				err := check(test)
				assert.NoErrorf(t, err, "LeadingDigit check should accept any digits passed to it, input was %s, got error: %s", test, err)
			}

			// Invalid leading digits
			for _, test := range []string{"0", "21", " 1", "-7"} {
				err := check(test)
				assert.Errorf(t, err, "LeadingDigit failed to error on invalid leading digit: %s", test)
			}
		}

	})
}

func TestDigitCountEquals(t *testing.T) {
	t.Run("Should correctly validate the number of digits in strings containing only digits", func(t *testing.T) {
		check := DigitCountEquals(5)

		if assert.NotNil(t, check) {

			// Test strings with 5 digits and nothing else in the string
			for _, test := range []string{"12345", "99999", "00000", "01010"} {
				err := check(test)
				assert.NoErrorf(t, err, "DigitCountEquals check failed to check the number of digits, input was %s, got error: %s", test, err)
			}

			// Test that the function generates an error when a string without 5 digits in it is given
			for _, test := range []string{"1", "12", "123", "1234", "123456", "0"} {
				err := check(test)
				assert.Errorf(t, err, "DigitCountEquals failed to error on an incorrect number of digits: %s", test)
			}
		}

	})

	t.Run("Should correctly validate the number of digits in strings containing a mix of chars", func(t *testing.T) {
		check := DigitCountEquals(5)

		if assert.NotNil(t, check) {

			// Test strings with 5 digits and other goodies in the string
			for _, test := range []string{"12-345", " 99999", "abcdef00000", "_01010"} {
				err := check(test)
				assert.NoErrorf(t, err, "DigitCountEquals check failed to check the number of digits, input was %s, got error: %s", test, err)
			}

			// Test that the function generates an error when a string without 5 digits in it is given
			for _, test := range []string{"1.", " 12", "-123a", " 1234", "...123456", "0-000"} {
				err := check(test)
				assert.Errorf(t, err, "DigitCountEquals failed to error on an incorrect number of digits: %s", test)
			}
		}

	})
}

func TestOnlyDigitsOrValidSeparator(t *testing.T) {
	t.Run("Should correctly validate digits and separator", func(t *testing.T) {
		check := OnlyDigitsOrValidSeparator("-", 4)
		for _, test := range []string{"123456", "4444-4444-4444-4444", "1234-5678"} {
			err := check(test)
			assert.NoErrorf(t, err, "OnlyDigitsOrValidSeparator should succeed for %s", test)
		}

		check = OnlyDigitsOrValidSeparator(" ", 2)
		for _, test := range []string{"123456", "44 44 44 44 44 44", "12 34"} {
			err := check(test)
			assert.NoErrorf(t, err, "OnlyDigitsOrValidSeparator should succeed for %s", test)
		}
	})

	t.Run("Should generate validation error for invalid input", func(t *testing.T) {
		check := OnlyDigitsOrValidSeparator("-", 4)
		for _, test := range []string{"O123456", "4444-4444-4444-4444-", "1234 5678 9090 1234", ""} {
			err := check(test)
			assert.Errorf(t, err, "OnlyDigitsOrValidSeparator should fail for %s", test)
		}

		check = OnlyDigitsOrValidSeparator(" ", 2)
		for _, test := range []string{"O123456", "4444-4444-4444-4444", "1234-5678", " 123"} {
			err := check(test)
			assert.Errorf(t, err, "OnlyDigitsOrValidSeparator should fail for %s", test)
		}
	})
}

func TestDoesNotExceedMaxRepeatingDigits(t *testing.T) {
	t.Run("Should correctly validate valid input data", func(t *testing.T) {
		check := DoesNotExceedMaxRepeatingDigits(3)
		for _, test := range []string{"1234567890", "1", "", "11223344", "112223344", "4121-1123-0001-4111", "4 1 1 4 4 1 4"} {
			err := check(test)
			assert.NoErrorf(t, err, "DoesNotExceedMaxRepeatingDigits should validate %s", test)
		}

		check = DoesNotExceedMaxRepeatingDigits(4)
		for _, test := range []string{"1234567890", "1", "", "11223344", "11112222344", "4121-1120-0001-4111", "4 1 4 4 4 4"} {
			err := check(test)
			assert.NoErrorf(t, err, "DoesNotExceedMaxRepeatingDigits should validate %s", test)
		}
	})

	t.Run("Should generate validation error for invalid input data", func(t *testing.T) {
		check := DoesNotExceedMaxRepeatingDigits(3)
		for _, test := range []string{"444400017654", "4444-1111-0000-1234", "4121-1110", "1 1 1 1"} {
			err := check(test)
			assert.Errorf(t, err, "DoesNotExceedMaxRepeatingDigits should fail for %s", test)
		}

		check = DoesNotExceedMaxRepeatingDigits(4)
		for _, test := range []string{"44444017654", "44444-1111-0000-1234", "4121-1111", "1 1 1 1 1"} {
			err := check(test)
			assert.Errorf(t, err, "DoesNotExceedMaxRepeatingDigits should fail for %s", test)
		}
	})
}
