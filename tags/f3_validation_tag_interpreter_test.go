package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Expected template:
// [{country_code}:{validation},{validation} | {country_code}:{validation}]
func Test_RetrieveValidateFieldsCreatesValidationStructAsExpected(t *testing.T) {
	cases := []struct {
		description                   string
		validationStr                 string
		hasErrors                     bool
		expectedErrorMessage          string
		expectedCountryValidationInfo map[string]*CountryValidationInfo
	}{
		{
			description:          "fails validation due to unexpected symbol [INITIAL_STATE]",
			validationStr:        ":GB:7-10,required | PT:5]",
			hasErrors:            true,
			expectedErrorMessage: "unexpected : symbol in position 0",
		},
		{
			description:          "fails validation due to unexpected symbol [ASSEMBLING_COUNTRY_CODE_STATE]",
			validationStr:        "[`GB:7-10,required | PT:5]",
			hasErrors:            true,
			expectedErrorMessage: "unexpected ` symbol in position 1",
		},
		{
			description:          "fails validation due to unexpected symbol [ASSEMBLING_COUNTRY_CODE_STATE]",
			validationStr:        "[GB->7-10,required | PT:5]",
			hasErrors:            true,
			expectedErrorMessage: "unexpected - symbol in position 3",
		},
		{
			description:          "fails validation due to unexpected symbol - after symbol :  [ASSEMBLING_COUNTRY_VALIDATION]",
			validationStr:        " [GB:-10 | PT:5]",
			hasErrors:            true,
			expectedErrorMessage: "unexpected - symbol in position 5",
		},
		{
			description:          "fails validation due to unexpected symbol + after symbol 1  [ASSEMBLING_COUNTRY_VALIDATION]",
			validationStr:        "[GB:1+10 | PT:5]",
			hasErrors:            true,
			expectedErrorMessage: "unexpected + symbol in position 5",
		},
		{
			description:                   "success validation with country GB and size equal to 1-10  [ASSEMBLING_COUNTRY_VALIDATION]",
			validationStr:                 " [GB:1-10]",
			hasErrors:                     false,
			expectedCountryValidationInfo: map[string]*CountryValidationInfo{"GB": {minLen: 1, maxLen: 10}},
		},
		{
			description:                   "success validation with country GB and size equal to 2-120  [ASSEMBLING_COUNTRY_VALIDATION]",
			validationStr:                 " [GB:12-120]",
			hasErrors:                     false,
			expectedCountryValidationInfo: map[string]*CountryValidationInfo{"GB": {minLen: 12, maxLen: 120}},
		},
		{
			description:                   "success validation with country GB and size equal to 10",
			validationStr:                 " [GB:10 ]",
			hasErrors:                     false,
			expectedCountryValidationInfo: map[string]*CountryValidationInfo{"GB": {minLen: 10, maxLen: 10}},
		},
		{
			description:                   "success validation with country GB and size equal to 10 to 765 and is required [ASSEMBLING_COUNTRY_VALIDATION]",
			validationStr:                 " [GB:10-765, required ]",
			hasErrors:                     false,
			expectedCountryValidationInfo: map[string]*CountryValidationInfo{"GB": {minLen: 10, maxLen: 765, required: true}},
		},
		{
			description:          "success validation with country GB and size equal to 10 to 765 and is required [ASSEMBLING_COUNTRY_VALIDATION]",
			validationStr:        " [GB:10-765, ]",
			hasErrors:            true,
			expectedErrorMessage: "unexpected ] symbol in position 13",
		},
		{
			description:          "success validation with country GB and size equal to 10 to 765 and is required [ASSEMBLING_COUNTRY_VALIDATION]",
			validationStr:        " [GB:10-765, | ]",
			hasErrors:            true,
			expectedErrorMessage: "unexpected | symbol in position 13",
		},
		{
			description:                   "success validation",
			validationStr:                 " [ GB:7-10,required | PT:5]",
			hasErrors:                     false,
			expectedCountryValidationInfo: map[string]*CountryValidationInfo{"GB": {minLen: 7, maxLen: 10, required: true}, "PT": {minLen: 5, maxLen: 5}},
		},
		{
			description:                   "success validation",
			validationStr:                 "[GB:7-10,required | PT:5 | AU:10-12, required]",
			hasErrors:                     false,
			expectedCountryValidationInfo: map[string]*CountryValidationInfo{"GB": {minLen: 7, maxLen: 10, required: true}, "PT": {minLen: 5, maxLen: 5}, "AU": {minLen: 10, maxLen: 12, required: true}},
		},
		{
			description:                   "success validation",
			validationStr:                 "[GB:required]",
			hasErrors:                     false,
			expectedCountryValidationInfo: map[string]*CountryValidationInfo{"GB": {required: true}},
		},
		{
			description:          "success validation",
			validationStr:        "[GB:7-10,required | AU:5 | AU:10-12, required]",
			hasErrors:            true,
			expectedErrorMessage: "country defined twice",
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			cInfo, err := CompileCountriesValidationInfos(c.validationStr)
			if c.hasErrors {
				assert.NotNil(t, err)
				assert.Equal(t, c.expectedErrorMessage, err.Error())
			} else {
				assert.Nil(t, err)
				assert.True(t, assertEquals(&c.expectedCountryValidationInfo, &cInfo))
			}

		})
	}
}

func Test_IsLetterWorks(t *testing.T) {
	cases := []struct {
		description    string
		char           byte
		expectedResult bool
	}{
		{description: "a is a letter", char: 'a', expectedResult: true},
		{description: "A is a letter", char: 'A', expectedResult: true},
		{description: "z is a letter", char: 'z', expectedResult: true},
		{description: "Z is a letter", char: 'Z', expectedResult: true},
		{description: "c is a letter", char: 'c', expectedResult: true},
		{description: "B is a letter", char: 'B', expectedResult: true},
		{description: "[ is NOT a letter", char: '[', expectedResult: false},
		{description: "- is NOT a letter", char: '-', expectedResult: false},
		{description: "= is NOT a letter", char: '=', expectedResult: false},
		{description: ", is NOT a letter", char: ',', expectedResult: false},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			result := IsLetter(c.char)
			assert.Equal(t, c.expectedResult, result)
		})
	}
}

func Test_IsNumericWorks(t *testing.T) {
	cases := []struct {
		description    string
		char           byte
		expectedResult bool
	}{
		{description: "1 is numeric", char: '1', expectedResult: true},
		{description: "2 is numeric", char: '2', expectedResult: true},
		{description: "3 is numeric", char: '3', expectedResult: true},
		{description: "4 is numeric", char: '4', expectedResult: true},
		{description: "5 is numeric", char: '5', expectedResult: true},
		{description: "6 is numeric", char: '6', expectedResult: true},
		{description: "7 is numeric", char: '7', expectedResult: true},
		{description: "8 is numeric", char: '8', expectedResult: true},
		{description: "9 is numeric", char: '9', expectedResult: true},
		{description: "[ is NOT numeric", char: '[', expectedResult: false},
		{description: "- is NOT numeric", char: '-', expectedResult: false},
		{description: "= is NOT numeric", char: '=', expectedResult: false},
		{description: ", is NOT numeric", char: ',', expectedResult: false},
		{description: "a is NOT numeric", char: 'a', expectedResult: false},
		{description: "A is NOT numeric", char: 'A', expectedResult: false},
		{description: "z is NOT numeric", char: 'z', expectedResult: false},
		{description: "Z is NOT numeric", char: 'Z', expectedResult: false},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			result := IsNumeric(c.char)
			assert.Equal(t, c.expectedResult, result)
		})
	}
}

func assertEquals(expectedCountryInfo, actualCountryInfo *map[string]*CountryValidationInfo) bool {
	for k, expected := range *expectedCountryInfo {
		actual := (*actualCountryInfo)[k]
		if actual == nil {
			return false
		}
		if !(expected.maxLen == actual.maxLen && expected.minLen == actual.minLen && expected.required == actual.required) {
			return false
		}
	}
	return true
}
