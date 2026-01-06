package xvalidator

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
	"github.com/shopspring/decimal"
)

// Decimal validation default values
const (
	// DefaultPrecision defines the default precision for decimal validation (total digits).
	DefaultPrecision = 38

	// DefaultScale defines the default scale for decimal validation (decimal places).
	DefaultScale = 18
)

// Decimal validation logic functions

// validateDecimalOperation creates a validator function for decimal operations.
// It handles string input and compares it using the provided comparator function.
func validateDecimalOperation(comparator func(d1, d2 *decimal.Decimal) bool) validator.Func {
	return func(fl validator.FieldLevel) bool {
		// Handle string input for decimal validation
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		// Parse field value as decimal
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}

		// Parse parameter value as decimal
		baseValue, err := decimal.NewFromString(fl.Param())
		if err != nil {
			return false
		}

		return comparator(&value, &baseValue)
	}
}

// Decimal comparison functions

// decimalGreaterThan compares if first decimal is greater than second.
func decimalGreaterThan(a, b *decimal.Decimal) bool {
	return a.GreaterThan(*b)
}

// decimalGreaterThanOrEqual compares if first decimal is greater than or equal to second.
func decimalGreaterThanOrEqual(a, b *decimal.Decimal) bool {
	return a.GreaterThanOrEqual(*b)
}

// decimalLessThan compares if first decimal is less than second.
func decimalLessThan(a, b *decimal.Decimal) bool {
	return a.LessThan(*b)
}

// decimalLessThanOrEqual compares if first decimal is less than or equal to second.
func decimalLessThanOrEqual(a, b *decimal.Decimal) bool {
	return a.LessThanOrEqual(*b)
}

// decimalEqual compares if two decimals are equal.
func decimalEqual(a, b *decimal.Decimal) bool {
	return a.Equal(*b)
}

// decimalNotEqual compares if two decimals are not equal.
func decimalNotEqual(a, b *decimal.Decimal) bool {
	return !a.Equal(*b)
}

// Mobile E.164 validation logic functions

// validateMobileE164 validates that the phone number is in E.164 format and is a mobile number.
// E.164 format: starts with +, followed by up to 15 digits.
// Uses libphonenumber to determine if the number is a mobile number.
// Supports country-specific validation:
//   - mobile_e164 (no param): validates any country mobile
//   - mobile_e164=TH: validates Thailand mobile numbers only
//   - mobile_e164=US: validates US mobile numbers only
//   - mobile_e164=XX: validates specific country mobile numbers
func validateMobileE164(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()

	// First check E.164 format with regex for performance
	if !E164Regex().MatchString(phoneNumber) {
		return false
	}

	// Parse the phone number without specifying region (let the library determine from prefix)
	num, err := phonenumbers.Parse(phoneNumber, "")
	if err != nil {
		return false
	}

	// Check if the number is valid
	if !phonenumbers.IsValidNumber(num) {
		return false
	}

	// Get the number type
	numberType := phonenumbers.GetNumberType(num)

	// Must be mobile type or fixed line or mobile (common in US and some countries)
	if numberType != phonenumbers.MOBILE && numberType != phonenumbers.FIXED_LINE_OR_MOBILE {
		return false
	}

	// Check country-specific validation if parameter is provided
	param := fl.Param()
	if param != "" {
		// Get the region code from the parsed number
		regionCode := phonenumbers.GetRegionCodeForNumber(num)

		// Compare with the expected country code
		if regionCode != param {
			return false
		}
	}
	return true
}

// URL validation logic functions

// validateHttpsScheme validates that the URL uses HTTPS scheme and has a valid host.
func validateHttpsScheme(fl validator.FieldLevel) bool {
	urlStr := fl.Field().String()
	parsed, err := url.Parse(urlStr)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return false
	}
	return true
}

// Decimal type registration function

// decimalTypeFunc returns the custom type function for decimal.Decimal registration.
func decimalTypeFunc(field reflect.Value) any {
	if valuer, ok := field.Interface().(decimal.Decimal); ok {
		return valuer.String()
	}
	return nil
}

// validateDecimal validates decimal precision and scale according to specified rules.
// Supports formats:
//   - decimal (default: precision=DefaultPrecision, scale=DefaultScale)
//   - decimal=2 (precision=DefaultPrecision, scale=2)
//   - decimal=0 (precision=DefaultPrecision, scale=0 - integer only)
//   - decimal=38:18 (precision=38, scale=18)
//   - decimal=10:6 (precision=10, scale=6)
func validateDecimal(fl validator.FieldLevel) bool {
	// Handle string input for decimal validation
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	// Parse field value as decimal
	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}

	// Parse parameters for precision and scale
	param := fl.Param()
	precision, scale := parseDecimalParams(param)

	// Validate precision and scale
	return validateDecimalPrecisionScale(value, precision, scale)
}

// parseDecimalParams parses decimal validation parameters.
// Returns precision and scale based on parameter format.
func parseDecimalParams(param string) (precision, scale int32) {
	// Default values
	precision = DefaultPrecision
	scale = DefaultScale

	if param == "" {
		return precision, scale
	}

	// Check if parameter contains colon (precision:scale format)
	if strings.Contains(param, ":") {
		parts := strings.Split(param, ":")
		if len(parts) == 2 {
			if p, err := strconv.ParseInt(parts[0], 10, 32); err == nil {
				precision = int32(p)
			}
			if s, err := strconv.ParseInt(parts[1], 10, 32); err == nil {
				scale = int32(s)
			}
		}
	} else {
		// Single parameter means scale only, precision defaults to DefaultPrecision
		if s, err := strconv.ParseInt(param, 10, 32); err == nil {
			scale = int32(s)
		}
	}

	return precision, scale
}

// validateDecimalPrecisionScale validates if decimal value fits within specified precision and scale.
func validateDecimalPrecisionScale(value decimal.Decimal, precision, scale int32) bool {
	// Get string representation of the decimal
	valueStr := value.String()

	// Handle negative numbers
	valueStr = strings.TrimPrefix(valueStr, "-")

	// Split by decimal point
	parts := strings.Split(valueStr, ".")
	integerPart := parts[0]
	decimalPart := ""

	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// Remove leading zeros from integer part (except if it's just "0")
	if integerPart != "0" {
		integerPart = strings.TrimLeft(integerPart, "0")
		if integerPart == "" {
			integerPart = "0"
		}
	}

	// Calculate integer digits and decimal places
	integerDigits := int32(len(integerPart))
	decimalPlaces := int32(len(decimalPart))

	// Validate scale (decimal places)
	if decimalPlaces > scale {
		return false
	}

	// Validate precision (integer digits + scale should not exceed precision)
	// For precision validation, we need to check if the integer part fits
	// within the available space after reserving space for the scale
	maxIntegerDigits := precision - scale
	return integerDigits <= maxIntegerDigits
}

// parseDecimalIfParam parses the decimal_if parameter.
// Parameter format: "rule@field=value"
// Examples:
//   - "2@Mode=mode1" -> rule="2", field="Mode", expect="mode1"
//   - "38:19@Mode=mode2" -> rule="38:19", field="Mode", expect="mode2"
//
// Returns rule (decimal format), field name, expected value, and error.
func parseDecimalIfParam(param string) (rule, field, expect string, err error) {
	// Split by @ to separate rule and condition
	parts := strings.Split(param, "@")
	if len(parts) != 2 {
		return "", "", "", validator.ValidationErrors{}
	}

	rule = parts[0]

	// Split condition by = to get field and expected value
	conditionParts := strings.Split(parts[1], "=")
	if len(conditionParts) != 2 {
		return "", "", "", validator.ValidationErrors{}
	}

	field = conditionParts[0]
	expect = conditionParts[1]

	return rule, field, expect, nil
}

// validateDecimalIf validates decimal precision and scale conditionally based on another field's value.
// Parameter format: "rule@field=value"
// Supports formats:
//   - decimal_if=2@Mode=mode1 -> if Mode equals "mode1", validate with scale 2 (precision=DefaultPrecision)
//   - decimal_if=38:19@Mode=mode2 -> if Mode equals "mode2", validate with precision 38 and scale 19
//   - decimal_if=0@Mode=mode3 -> if Mode equals "mode3", validate with scale 0 (integer only)
//   - decimal_if=@Mode=mode4 -> if Mode equals "mode4", use default precision and scale
func validateDecimalIf(fl validator.FieldLevel) bool {
	rule, field, expect, err := parseDecimalIfParam(fl.Param())
	if err != nil {
		return false
	}

	// Read other field value to check condition
	parent := fl.Parent()
	otherField := parent.FieldByName(field)
	if !otherField.IsValid() {
		return false
	}

	other := otherField.String()
	if other != expect {
		return true // Condition not met → skip validation
	}

	// Handle string input for decimal validation (same as validateDecimal)
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	// Parse field value as decimal
	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}

	// Parse parameters for precision and scale using same logic as decimal rule
	precision, scale := parseDecimalParams(rule)

	// Validate precision and scale using same logic as decimal rule
	return validateDecimalPrecisionScale(value, precision, scale)
}

// Password validation logic functions

// validatePasswordStrength validates password strength according to security requirements.
// Password must meet the following criteria:
//   - At least 8 characters long
//   - Contains at least one uppercase letter (A-Z)
//   - Contains at least one lowercase letter (a-z)
//   - Contains at least one digit (0-9)
//   - Contains at least one special character (!@#$%^&*()_+-=[]{}|;:,.<>?)
func validatePasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if err := ValidatePasswordStrength(password); err != nil {
		return false
	}

	return true
}
