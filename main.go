package main

import (
	"reflect"

	"github.com/GUILN/custom-validation-tags/tags"
)

type ValidationMatrix map[string]*map[string]*tags.CountryValidationInfo

func main() {
	// acc := &tags.Account{
	// 	Country: "GB",
	// 	BankId:  "asdfasdfasdfasfasdfasdf",
	// }

	//bankIdField, _ := reflect.TypeOf(*acc).FieldByName("BankId")
	//t := bankIdField.Tag.Get("f3_validate")

	// fmt.Printf("Account: \n%v", acc)
	// fmt.Printf("Form3 validation Tag: \n%s", t)
	// validationMatrix := CreateValidationMatrix(*acc)
	// for _, field := range validationMatrix {
	// 	for _, countryValidationInfo := range *field {
	// 		fmt.Printf("%v\n", countryValidationInfo)
	// 	}
	// }
	//fmt.Println(validationMatrix)
}

// func validate(account *tags.Account) {
// 	t := reflect.TypeOf(*account)

// 	for i := 0; i < t.NumField(); i++ {
// 		fmt.Printf("Field: %s\n", t.Field(i).Name)
// 	}
// }

func CreateValidationMatrix(i interface{}) ValidationMatrix {
	t := reflect.TypeOf(i)

	matrix := make(map[string]*map[string]*tags.CountryValidationInfo)

	for i := 0; i < t.NumField(); i++ {
		validationTag := t.Field(i).Tag.Get(tags.ValidationForm3TagName)
		if len(validationTag) == 0 {
			break
		}
		validationInfos, _ := tags.CompileCountriesValidationInfos(validationTag)
		fieldName := t.Field(i).Name

		matrix[fieldName] = &validationInfos
	}

	return matrix
}
