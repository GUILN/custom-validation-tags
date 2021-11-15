package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type account struct {
	Country string
	BankId  string `f3_validate:"[GB:7-10,required | PT:5]"`
	IBAN    string `f3_validate:"[GB:8 | AU:4,required | PT:7-9,required]"`
}

func Test_CreateValidationMatrix_CreatesExpectedValidationMatrix(t *testing.T) {
	acc := &account{
		Country: "GB",
		BankId:  "123344",
		IBAN:    "asdf",
	}

	validationMatrix, err := CreateValidationMatrix(*acc)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(validationMatrix))
	assert.Equal(t, 4, ((*validationMatrix["IBAN"])["AU"]).maxLen)
}
