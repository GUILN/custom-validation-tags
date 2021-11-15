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

func Test_GivenTestCase_WhenICallValidateMethod_ThenItReturnsExpectedResult(t *testing.T) {

	cases := []struct {
		description              string
		acc                      *account
		expectedValidationErrors []string
	}{
		{
			description:              "when account is invalid then validation returns expected validation errors",
			acc:                      &account{Country: "GB", BankId: "123"},
			expectedValidationErrors: []string{"field BankId must have size from 7 to 10 when country is GB but found size 3"},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			validationResult, err := Validate(*c.acc, c.acc.Country)
			assert.Nil(t, err)

			if c.expectedValidationErrors == nil {
				assert.Nil(t, validationResult)
			} else {
				for _, expectedValidationErr := range c.expectedValidationErrors {
					assert.Contains(t, validationResult, expectedValidationErr)
				}
			}

		})
	}
}

type wrongAccountStruct struct {
	Country string
	BankId  string `f3_validate:"[GB|]"`
}

func Test_WhenICallValidateMethodWithTagsInWrongFormat_ThenIGetABuildError(t *testing.T) {
	acc := &wrongAccountStruct{
		Country: "",
		BankId:  "",
	}
	_, err := Validate(*acc, acc.Country)

	assert.NotNil(t, err)
}
