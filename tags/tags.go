package tags

type Account struct {
	Country string
	BankId  string `f3_validate:"[GB:7-10,required | PT:5]"`
}
