package xvalidator

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

// RegisterDecimalValidators registers all decimal validation rules for shopspring/decimal.Decimal type.
// This function adds comprehensive decimal comparison validators that work with decimal.Decimal values.
func RegisterDecimalValidators(v *validator.Validate) {
	// Register decimal comparison operations
	v.RegisterValidation("dgt", validateDecimalOperation(decimalGreaterThan))
	v.RegisterValidation("dgte", validateDecimalOperation(decimalGreaterThanOrEqual))
	v.RegisterValidation("dlt", validateDecimalOperation(decimalLessThan))
	v.RegisterValidation("dlte", validateDecimalOperation(decimalLessThanOrEqual))
	v.RegisterValidation("deq", validateDecimalOperation(decimalEqual))
	v.RegisterValidation("dneq", validateDecimalOperation(decimalNotEqual))

	// Register decimal precision and scale validation
	v.RegisterValidation("decimal", validateDecimal)

	// Register conditional decimal validation
	v.RegisterValidation("decimal_if", validateDecimalIf)

	// Register decimal type for proper handling
	v.RegisterCustomTypeFunc(decimalTypeFunc, decimal.Decimal{})
}

// RegisterURLValidators registers URL-specific validation rules.
// This function adds validators for URL format and protocol validation.
func RegisterURLValidators(v *validator.Validate) {
	v.RegisterValidation("https_url", validateHttpsScheme)
}

// RegisterPhoneValidators registers phone number validation rules using libphonenumber.
// This function adds validators for international phone number format and type validation.
func RegisterPhoneValidators(v *validator.Validate) {
	v.RegisterValidation("mobile_e164", validateMobileE164)
}

// RegisterPasswordValidators registers password validation rules.
// This function adds validators for password strength and complexity requirements.
func RegisterPasswordValidators(v *validator.Validate) {
	v.RegisterValidation("password_strength", validatePasswordStrength)
}
