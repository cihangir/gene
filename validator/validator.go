package validator

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// ERR_MSG_X_MUST_BE_OF_TYPE_Y = `%s must be of type %s`

	// ERR_MSG_X_IS_MISSING_AND_REQUIRED  = `%s is missing and required`
	// ERR_MSG_MUST_BE_OF_TYPE_X          = `must be of type %s`
	// ERR_MSG_ARRAY_ITEMS_MUST_BE_UNIQUE = `array items must be unique`

	ERR_MSG_STRING_LENGTH_MUST_BE_GREATER_OR_EQUAL = `string length must be greater or equal to %d`
	ERR_MSG_STRING_LENGTH_MUST_BE_LOWER_OR_EQUAL   = `string length must be lower or equal to %d`
	ERR_MSG_DOES_NOT_MATCH_PATTERN                 = `does not match pattern '%s'`
	ERR_MSG_MUST_MATCH_ONE_ENUM_VALUES             = `must match one of the enum values [%s]`

	ERR_MSG_NUMBER_MUST_BE_GREATER = `must be greater than %f`
	ERR_MSG_NUMBER_MUST_BE_LOWER   = `must be lower than %f`
	ERR_MSG_MULTIPLE_OF            = `must be a multiple of %f`

	// ERR_MSG_NUMBER_MUST_BE_LOWER_OR_EQUAL   = `must be lower than or equal to %s`
	// ERR_MSG_NUMBER_MUST_BE_GREATER_OR_EQUAL = `must be greater than or equal to %f`

	// ERR_MSG_NUMBER_MUST_VALIDATE_ALLOF = `must validate all the schemas (allOf)`
	// ERR_MSG_NUMBER_MUST_VALIDATE_ONEOF = `must validate one and only one schema (oneOf)`
	// ERR_MSG_NUMBER_MUST_VALIDATE_ANYOF = `must validate at least one schema (anyOf)`
	// ERR_MSG_NUMBER_MUST_VALIDATE_NOT   = `must not validate the schema (not)`

	// ERR_MSG_ARRAY_MIN_ITEMS = `array must have at least %d items`
	// ERR_MSG_ARRAY_MAX_ITEMS = `array must have at the most %d items`

	// ERR_MSG_ARRAY_MIN_PROPERTIES = `must have at least %d properties`
	// ERR_MSG_ARRAY_MAX_PROPERTIES = `must have at the most %d properties`

	// ERR_MSG_HAS_DEPENDENCY_ON = `has a dependency on %s`

	// ERR_MSG_ARRAY_NO_ADDITIONAL_ITEM = `no additional item allowed on array`

	// ERR_MSG_ADDITIONAL_PROPERTY_NOT_ALLOWED = `additional property "%s" is not allowed`
	// ERR_MSG_INVALID_PATTERN_PROPERTY        = `property "%s" does not match pattern %s`
)

func MinLength(data string, length int) error {
	if utf8.RuneCount([]byte(data)) < length {
		return fmt.Errorf(ERR_MSG_STRING_LENGTH_MUST_BE_GREATER_OR_EQUAL, length)
	}

	return nil
}

func MaxLength(data string, length int) error {
	if utf8.RuneCount([]byte(data)) > length {
		return fmt.Errorf(ERR_MSG_STRING_LENGTH_MUST_BE_LOWER_OR_EQUAL, length)
	}

	return nil
}

func Pattern(data string, pattern string) error {
	// TODO add caching for compile?
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	if !regex.MatchString(data) {
		return fmt.Errorf(ERR_MSG_DOES_NOT_MATCH_PATTERN, pattern)
	}

	return nil
}

func OneOf(data string, enums []string) error {
	for _, val := range enums {
		if val == data {
			return nil
		}
	}

	return fmt.Errorf(ERR_MSG_MUST_MATCH_ONE_ENUM_VALUES, strings.Join(enums, ","))
}

func Min(data float64, min float64) error {
	if data < min {
		return fmt.Errorf(ERR_MSG_NUMBER_MUST_BE_GREATER, min)
	}

	return nil
}

func Max(data float64, max float64) error {
	if data > max {
		return fmt.Errorf(ERR_MSG_NUMBER_MUST_BE_LOWER, max)
	}

	return nil
}

func MultipleOf(data float64, multipleOf float64) error {
	if math.Mod(data, multipleOf) == 0 {
		return fmt.Errorf(ERR_MSG_MULTIPLE_OF, multipleOf)
	}

	return nil
}