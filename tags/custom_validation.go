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
	validationCloser             Symbol = ']'
	countryValidationInitializer Symbol = ':'

	countrySeparator       Symbol = '|'
	numericLengthSeparator Symbol = '-'
	validationSeparator    Symbol = ','
)

// Tokens

var tokenMap = map[string]func(*CountryValidationInfo){
	"required": func(countryValidationInfo *CountryValidationInfo) {
		countryValidationInfo.required = true
	},
}

// States
type State string

const (
	invalidState                             State = "[INVALID]"
	initialState                             State = "[INITIAL]"
	assemblingCountryCode                    State = "[ASSB_COUNTRY_CODE]"
	assemblingCountryValidation              State = "[ASSB_COUNTRY_VALIDATION]"
	assemblingContryValidationFieldSize      State = "[ASSB_COUNTRY_VALIDATION_FLD_SIZE]"
	assemblingContryValidationFieldSizeMax   State = "[ASSB_COUNTRY_VALIDATION_FLD_SIZE_MAX]"
	assemblingCountryValidationToken         State = "[ASSB_COUNTRY_VALIDATION_FLD_TOKEN]"
	expectingCountryValidationCloseStatement State = "[ASSB_COUNTRY_VALIDATION_EXPECTING_CLOSE_STATEMENT]"
	finalState                               State = "[FINAL]"
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
		} else if entrySymbol == ' ' && *currentCountry == "" {
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
			*accumulator += string(entrySymbol)
			return assemblingCountryValidationToken, nil
		} else if IsNumeric(entrySymbol) {
			*accumulator += string(entrySymbol)
			return assemblingContryValidationFieldSize, nil
		} else if entrySymbol == ' ' {
			return assemblingCountryValidation, nil
		}
		return invalidState, createUnexpectedSymbolError(entrySymbol, position)
	},
	expectingCountryValidationCloseStatement: func(entrySymbol byte, countries *map[string]*CountryValidationInfo, accumulator, currentCountry *string, position int) (State, error) {
		if entrySymbol == byte(validationCloser) {
			return finalState, nil
		} else if entrySymbol == byte(validationSeparator) {
			return assemblingCountryValidation, nil
		} else if entrySymbol == ' ' {
			return expectingCountryValidationCloseStatement, nil
		} else if IsLetter(entrySymbol) {
			*accumulator += string(entrySymbol)
			return assemblingCountryValidationToken, nil
		} else if entrySymbol == byte(countrySeparator) {
			*currentCountry = ""
			return assemblingCountryCode, nil
		}
		return invalidState, createUnexpectedSymbolError(entrySymbol, position)
	},
	assemblingContryValidationFieldSize: func(entrySymbol byte, countries *map[string]*CountryValidationInfo, accumulator *string, currentCountry *string, position int) (State, error) {
		if IsNumeric(entrySymbol) {
			*accumulator += string(entrySymbol)
			return assemblingContryValidationFieldSize, nil
		} else if entrySymbol == byte(numericLengthSeparator) {
			(*countries)[*currentCountry].minLen = ParseToInt(*accumulator)
			*accumulator = ""
			return assemblingContryValidationFieldSizeMax, nil
		} else if entrySymbol == byte(validationCloser) {
			(*countries)[*currentCountry].maxLen = ParseToInt(*accumulator)
			(*countries)[*currentCountry].minLen = ParseToInt(*accumulator)
			*accumulator = ""
			return finalState, nil
		} else if entrySymbol == byte(validationSeparator) {
			(*countries)[*currentCountry].maxLen = ParseToInt(*accumulator)
			(*countries)[*currentCountry].minLen = ParseToInt(*accumulator)
			*accumulator = ""
			return assemblingCountryValidation, nil
		} else if entrySymbol == ' ' {
			(*countries)[*currentCountry].maxLen = ParseToInt(*accumulator)
			(*countries)[*currentCountry].minLen = ParseToInt(*accumulator)
			*accumulator = ""
			return expectingCountryValidationCloseStatement, nil
		}
		return invalidState, createUnexpectedSymbolError(entrySymbol, position)
	},
	assemblingContryValidationFieldSizeMax: func(entrySymbol byte, countries *map[string]*CountryValidationInfo, accumulator *string, currentCountry *string, position int) (State, error) {
		if IsNumeric(entrySymbol) {
			*accumulator += string(entrySymbol)
			return assemblingContryValidationFieldSizeMax, nil
		} else if entrySymbol == byte(validationSeparator) {
			(*countries)[*currentCountry].maxLen = ParseToInt(*accumulator)
			*accumulator = ""
			return assemblingCountryValidation, nil
		} else if entrySymbol == byte(validationCloser) {
			(*countries)[*currentCountry].maxLen = ParseToInt(*accumulator)
			*accumulator = ""
			return finalState, nil
		} else if entrySymbol == ' ' {
			(*countries)[*currentCountry].maxLen = ParseToInt(*accumulator)
			*accumulator = ""
			return expectingCountryValidationCloseStatement, nil
		}
		return invalidState, createUnexpectedSymbolError(entrySymbol, position)
	},
	assemblingCountryValidationToken: func(entrySymbol byte, countries *map[string]*CountryValidationInfo, accumulator *string, currentCountry *string, position int) (State, error) {
		if IsLetter(entrySymbol) {
			*accumulator += string(entrySymbol)
			return assemblingCountryValidationToken, nil
		} else if entrySymbol == ' ' || entrySymbol == byte(validationCloser) {
			tokenFunction := tokenMap[*accumulator]
			if tokenFunction == nil {
				return invalidState, createUnexpectedTokenError(*accumulator, position)
			}
			tokenFunction((*countries)[*currentCountry])
			*accumulator = ""

			if entrySymbol == ' ' {
				return expectingCountryValidationCloseStatement, nil
			}
			return finalState, nil
		}
		return assemblingCountryValidationToken, nil
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

		if currentState == finalState {
			return countriesValidationInfos, nil
		}
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

func ParseToInt(str string) int {
	parsedInt, _ := strconv.ParseFloat(str, 64)
	return int(parsedInt)
}

func createUnexpectedSymbolError(unexpectedSymbol byte, position int) error {
	return fmt.Errorf("unexpected %s symbol in position %d", string(unexpectedSymbol), position)
}

func createUnexpectedTokenError(unexpectedToken string, position int) error {
	return fmt.Errorf("unexpected token %s in position %d", unexpectedToken, position)
}
