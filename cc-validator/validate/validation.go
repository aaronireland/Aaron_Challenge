package validate

// Check function validates the input for a specific condition
type Check func(input string) error

// Validation struct stores an array of functions which will
// ensure that the input provided meets all the conditions
type Validation struct {
	required []Check
}

// Require function addds the Check function in with the correct scope
func (v *Validation) Require(check Check) {
	v.required = append(v.required, check)
}

// Validate function validates the input provided for all the conditions
func (v *Validation) Validate(input string) (validated bool, errs []error) {
	for _, check := range v.required {
		if err := check(input); err != nil {
			errs = append(errs, err)
		}
	}

	validated = len(errs) == 0

	return validated, errs
}
