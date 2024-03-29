package datatype

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gobuffalo/nulls"
)

const (
	// regular expression to validation rzp standard ID
	RegexID = `^[a-zA-Z0-9]{14}$`

	// regular expression to validation unix timestamp
	RegexUnixTimestamp = `^([\d]{10}|0)$`

	// regular expression to validation basic string
	// consisting alpha, number, _, - and white space
	RegexBasicString = `^[a-zA-Z0-9_\-\s]*$`
)

// ValidateNullableInt64 checks NullableInt64 against the rules provided
func ValidateNullableInt64(value nulls.Int64, rules ...validation.RuleFunc) validation.RuleFunc {
	return func(value interface{}) error {
		v := value.(nulls.Int64)
		for _, f := range rules {
			if err := f(v.Int64); err != nil {
				return err
			}
		}

		return nil
	}
}

// ValidateNullableString checks NullableString against the rules provided
func ValidateNullableString(value nulls.String, rules ...validation.RuleFunc) validation.RuleFunc {
	return func(value interface{}) error {
		v := value.(nulls.String)
		for _, f := range rules {
			if err := f(v.String); err != nil {
				return err
			}
		}

		return nil
	}
}

// IsUUID checks the given string is a valid 14-char Unique ID
func IsUUID(value interface{}) error {
	return isValidString(value, RegexID)
}

// IsTimestamp will validate if the value is a valid unix timestamp or not
func IsTimestamp(value interface{}) error {
	if value == nil {
		return nil
	}

	return MatchRegex(fmt.Sprintf("%v", value), RegexUnixTimestamp)
}

// MatchRegex checks if given input matches a given regex or not
func MatchRegex(value string, regex string) error {
	if validString, err := regexp.Compile(regex); err != nil {
		return errors.New("invalid regex")
	} else if !validString.MatchString(value) {
		return errors.New("not a valid input")
	}

	return nil
}

// isValidString checks if given input matches a given regex or not
func isValidString(value interface{}, regex string) error {
	// let the nil handled by required validation
	if value == nil {
		return nil
	}

	if str, err := isString(value); err != nil {
		return err
	} else if str == "" {
		return nil
	} else {
		return MatchRegex(str, regex)
	}
}

// isString checks if the given data is valid string or not
func isString(value interface{}) (string, error) {
	if str, ok := value.(string); !ok {
		return "", errors.New("must be a string")
	} else {
		return str, nil
	}
}

// IsBasicString validates if the given value is basic string
// consisting alphabet, number, _, - and white space
func IsBasicString(value interface{}) error {
	return isValidString(value, RegexBasicString)
}

// IsInt64 checks if the given data is valid int64 or not
func IsInt64(value interface{}) error {
	if _, ok := value.(int64); !ok {
		return errors.New("must be an int64 integer")
	}
	return nil
}

// IsJson validates if the given value is valid json
func IsJson(value interface{}) error {
	if value == nil {
		return nil
	}

	if str, err := isString(value); err != nil {
		return err
	} else {
		input := []byte(str)
		var x struct{}
		if err := json.Unmarshal(input, &x); err != nil {
			return err
		}
	}

	return nil
}

// IsNumeric checks if the given string 's' is float, int, signed / unsigned, exponential
// and returns true if it is valid ..
func IsNumeric(value interface{}) error {
	valueStr, err := isString(value)
	if err != nil {
		return err
	}
	_, err = strconv.ParseFloat(valueStr, 64)

	return err
}

// IsBool checks if the given data is boolean or not
func IsBool(value interface{}) error {
	switch value.(type) {
	case nil:
		return nil
	case bool:
		return nil
	default:
		return errors.New("must be a boolean")
	}
}
