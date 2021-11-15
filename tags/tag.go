package tags

import "reflect"

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
