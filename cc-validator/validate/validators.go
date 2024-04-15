// validators.go contains all the templated validation checks. These return wrapped functions to allow
// configuration and because compiling regex patterns are expensive, it's more performant to compile them
// one time and then run the check as many times as neeeded.
package validate

import (
	"fmt"
	"regexp"
	"strconv"
)

// LeadingDigit function validates that the input starts with one of the specified numbers,
// if no digits are provided to the function all digits 0 to 9 are valid
func LeadingDigit(i ...int) Check {
	regexStr := "^"
	if len(i) > 0 {
		regexStr += "["
		for _, digit := range i {
			regexStr = regexStr + strconv.Itoa(digit)
		}
		regexStr += "]"
	} else {
		regexStr = `^\d`
	}

	r := regexp.MustCompile(regexStr)

	return func(input string) error {
		if r.MatchString(input) {
			return nil
		}
		return fmt.Errorf("invalid leading digit")
	}
}

// DigitCountEquals function validates that the input contains the correct number of digits
func DigitCountEquals(i int) Check {
	r := regexp.MustCompile(`\D`)

	return func(input string) error {
		digits := r.ReplaceAllString(input, "")
		if len(digits) != i {
			return fmt.Errorf("must contain exactly %d numeric digits", i)
		}
		return nil
	}
}

// OnlyDigitsOrValidSeparator validates that the input only contains numeric digits or
// that the input contains a sequence of digits with a validate separator in between
func OnlyDigitsOrValidSeparator(sep string, seq int) Check {
	regexStr := fmt.Sprintf(`^\d+$|^(?:\d{%v}|%s)+(\d{%v})$`, seq, sep, seq)
	r := regexp.MustCompile(regexStr)

	return func(input string) error {
		if r.MatchString(input) {
			return nil
		}
		return fmt.Errorf("only digits with an optional separator(%s) between blocks of %v are valid", sep, seq)
	}
}

// DoesNotExceedMaxRepeatingDigits function validates that no single numeric digit is repeated more than
// the maximum specified times
func DoesNotExceedMaxRepeatingDigits(max int) Check {
	var (
		count    = 0
		previous byte
	)
	return func(input string) error {
		for _, current := range []byte(input) {
			if current >= 48 && current <= 57 {
				switch previous {
				case 0:
					count++
					previous = current
				case current:
					count++
					if count > max {
						return fmt.Errorf("invalid digit sequence: %s occurred %v times in a row. max allowable is %v", string(current), count, max)
					}
				default:
					previous = current
					count = 1
				}
			}
		}
		return nil
	}
}
