package main

import (
	"fmt"
	"reflect"

	"sandbox.io/tags/tags"
)

func main() {
	acc := &tags.Account{
		Country: "GB",
		BankId:  "asdfasdfasdfasfasdfasdf",
	}

	bankIdField, _ := reflect.TypeOf(*acc).FieldByName("BankId")
	t := bankIdField.Tag.Get("f3_validate")

	fmt.Printf("Account: \n%v", acc)
	fmt.Printf("Form3 validation Tag: \n%s", t)
}

func validateAccount(acc *tags.Account) bool {
	return true
}
