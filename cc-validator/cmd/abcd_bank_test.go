package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AbcdBankValidator(t *testing.T) {
	gotValidator := NewAbcdBankValidator()
	assert.NotNil(t, gotValidator, "the AbcdBankValidator constructor returned a nil validator")
}
