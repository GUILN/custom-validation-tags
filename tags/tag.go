package tags

import (
	"fmt"
	"reflect"
)

const ValidationForm3TagName string = "f3_validate"

type ValidationMatrix map[string]*map[string]*CountryValidationInfo

func CreateValidationMatrix(i interface{}) (ValidationMatrix, error) {
	t := reflect.TypeOf(i)

	matrix := make(map[string]*map[string]*CountryValidationInfo)

	for i := 0; i < t.NumField(); i++ {
		validationTag := t.Field(i).Tag.Get(ValidationForm3TagName)
		if len(validationTag) > 0 {
			validationInfos, err := CompileCountriesValidationInfos(validationTag)
			if err != nil {
				return nil, err
			}
			fieldName := t.Field(i).Name

			matrix[fieldName] = &validationInfos
		}
	}

	return matrix, nil
}

func Validate(i interface{}, country string) ([]string, error) {
	validationMatrix, err := CreateValidationMatrix(i)
	var validationErrors []string = nil

	if err != nil {
		return nil, err
	}

	for fieldName, validationCountryMap := range validationMatrix {
		fieldValue := getFieldValueByFieldName(i, fieldName)
		if countryValidationInfo := (*validationCountryMap)[country]; countryValidationInfo != nil {
			if errs := getValidationErrors(country, fieldName, fieldValue, countryValidationInfo); errs != nil {
				validationErrors = append(validationErrors, errs...)
			}
		}
	}

	return validationErrors, nil
}

func getFieldValueByFieldName(i interface{}, fieldName string) string {
	r := reflect.ValueOf(i)
	f := reflect.Indirect(r).FieldByName(fieldName)
	return f.String()
}

func getValidationErrors(country, fieldName, fieldValue string, validationInfo *CountryValidationInfo) []string {
	var validationErrors []string = nil
	if validationInfo.minLen > 0 && validationInfo.maxLen > 0 {
		if validationInfo.minLen != validationInfo.maxLen {
			actualLen := len(fieldValue)
			if actualLen < validationInfo.minLen || actualLen > validationInfo.maxLen {
				validationErrors = append(validationErrors, fmt.Sprintf("field %s must have size from %d to %d when country is GB but found size %d", fieldName, validationInfo.minLen, validationInfo.maxLen, actualLen))
			}
		}
	}

	return validationErrors
}
