package tags

import (
	"fmt"
	"strconv"
)

type CountryValidationInfo struct {
	fieldName string
	country   string
	minLen    int
	maxLen    int
	required  bool
}

//Symbols
type Symbol byte

const (
	validationOpener             Symbol = '['
	validationCloser             Symbol = '}'
	countryValidationInitializer Symbol = ':'

	countrySeparator       Symbol = '|'
	numericLengthSeparator Symbol = '-'
	validationSeparator    Symbol = ','
)

// Alphabet

var Alphabet = []string{
	string(validationOpener),
	string(validationCloser),
	string(countryValidationInitializer),
	string(countrySeparator),
	string(numericLengthSeparator),
	string(validationSeparator),
	"required",
}

// States
type State string

const (
	invalidState                           State = "[INVALID]"
	initialState                           State = "[INITIAL]"
	assemblingCountryCode                  State = "[ASSB_COUNTRY_CODE]"
	assemblingCountryValidation            State = "[ASSB_COUNTRY_VALIDATION]"
	assemblingContryValidationFieldSize    State = "[ASSB_COUNTRY_VALIDATION_FLD_SIZE]"
	assemblingContryValidationFieldSizeMax State = "[ASSB_COUNTRY_VALIDATION_FLD_SIZE_MAX]"
	finalState                             State = "[FINAL]"
)

type TransictionFunction func(byte, *map[string]*CountryValidationInfo, *string, *string, int) (State, error)

var transictionTable = map[State]TransictionFunction{

	initialState: func(entrySymbol byte, countries *map[string]*CountryValidationInfo, currentCountry *string, accumulator *string, position int) (State, error) {
		if entrySymbol == byte(validationOpener) {
			return assemblingCountryCode, nil
		} else if entrySymbol == ' ' {
			return initialState, nil
		}
		return invalidState, createUnexpectedSymbolError(entrySymbol, position)
	},
	assemblingCountryCode: func(entrySymbol byte, countries *map[string]*CountryValidationInfo, accumulator *string, currentCountry *string, position int) (State, error) {
		if IsLetter(entrySymbol) {
			*currentCountry += string(entrySymbol)

			return assemblingCountryCode, nil
		} else if entrySymbol == byte(countryValidationInitializer) {
			if _, alreadyContaisCountry := (*countries)[*currentCountry]; alreadyContaisCountry {
				return invalidState, fmt.Errorf("country already exists")
			}

			(*countries)[*currentCountry] = &CountryValidationInfo{}
			return assemblingCountryValidation, nil
		}
		return invalidState, createUnexpectedSymbolError(entrySymbol, position)
	},
	assemblingCountryValidation: func(entrySymbol byte, countries *map[string]*CountryValidationInfo, accumulator *string, currentCountry *string, position int) (State, error) {
		if IsLetter(entrySymbol) {

		} else if IsNumeric(entrySymbol) {
			*accumulator += string(entrySymbol)
			return assemblingContryValidationFieldSize, nil
		}
		return invalidState, createUnexpectedSymbolError(entrySymbol, position)
	},
	assemblingContryValidationFieldSize: func(entrySymbol byte, countries *map[string]*CountryValidationInfo, accumulator *string, currentCountry *string, position int) (State, error) {
		if IsNumeric(entrySymbol) {
			*accumulator += string(entrySymbol)
			return assemblingContryValidationFieldSize, nil
		} else if entrySymbol == byte(numericLengthSeparator) {
			(*countries)[*currentCountry].minLen = 0
			*accumulator = ""
			return assemblingContryValidationFieldSizeMax, nil
		} else if entrySymbol == byte(validationSeparator) {
			(*countries)[*currentCountry].maxLen = 0
			(*countries)[*currentCountry].minLen = 0
			*accumulator = ""
			return assemblingCountryValidation, nil
		}
		return invalidState, createUnexpectedSymbolError(entrySymbol, position)
	},
}

func mountCountriesValidationInfos(validationStr string) (map[string]*CountryValidationInfo, error) {
	currentState := initialState
	var currentSymbol byte
	currentCountry := new(string)
	accumulator := new(string)
	var stateError error = nil
	countriesValidationInfos := make(map[string]*CountryValidationInfo)

	for i := 0; i < len(validationStr); i++ {
		currentSymbol = validationStr[i]
		currentState, stateError = transictionTable[currentState](currentSymbol, &countriesValidationInfos, accumulator, currentCountry, i)

		if currentState == invalidState {
			return nil, stateError
		}
	}

	return countriesValidationInfos, nil
}

func IsLetter(c byte) bool {
	return !((c < 'a' || c > 'z') && (c < 'A' || c > 'Z'))
}

func IsNumeric(c byte) bool {
	_, err := strconv.ParseFloat(string(c), 64)
	return err == nil
}

func createUnexpectedSymbolError(unexpectedSymbol byte, position int) error {
	return fmt.Errorf("unexpected %s symbol in position %d", string(unexpectedSymbol), position)
}
