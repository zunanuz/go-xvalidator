package xvalidator

import (
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator with enhanced decimal support and Universal Translator.
type Validator struct {
	validate   *validator.Validate
	translator ut.Translator
}

// NewValidator creates a new validator instance with all custom rules and English translator registered.
func NewValidator() (*Validator, error) {
	v := validator.New()

	// Register JSON tag name function for better field naming
	v.RegisterTagNameFunc(getJSONTagName)

	// Register all custom validators
	RegisterDecimalValidators(v)
	RegisterURLValidators(v)
	RegisterPhoneValidators(v)
	RegisterPasswordValidators(v)

	// Setup English translator
	trans, err := setupTranslator(v)
	if err != nil {
		return nil, err
	}

	return &Validator{
		validate:   v,
		translator: trans,
	}, nil
}

// GetTranslator returns the Universal Translator instance.
func (v *Validator) GetTranslator() ut.Translator {
	return v.translator
}

// GetValidator returns the underlying validator.Validate instance.
func (v *Validator) GetValidator() *validator.Validate {
	return v.validate
}

// Validate validates a struct and returns raw validation errors without translation.
// For user-friendly error messages, use StructTranslated instead.
func (v *Validator) Validate(i any) error {
	return v.validate.Struct(i)
}

// Struct validates a struct and returns raw validation errors without translation.
// This method is an alias for Validate for consistency with other validator methods.
func (v *Validator) Struct(i any) error {
	return v.validate.Struct(i)
}

// Var validates a single variable using the provided validation tag and returns raw errors.
// For user-friendly error messages, use VarTranslated instead.
func (v *Validator) Var(field any, tag string) error {
	return v.validate.Var(field, tag)
}

// StructTranslated validates a struct based on tags and returns user-friendly translated error messages.
func (v *Validator) StructTranslated(s any) error {
	err := v.validate.Struct(s)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatTranslatedErrors(validationErrors, v.translator)
		}
	}
	return err
}

// VarTranslated validates a single variable using the provided validation tag and returns user-friendly translated error messages.
func (v *Validator) VarTranslated(field any, tag string) error {
	err := v.validate.Var(field, tag)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatTranslatedErrors(validationErrors, v.translator)
		}
	}
	return err
}

// getJSONTagName extracts the JSON field name from a struct field's json tag.
// It handles cases where the tag contains options like "omitempty" or "-".
// Returns the field name if no json tag is present.
// Optimized version using strings.IndexByte for better performance.
func getJSONTagName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return field.Name
	}

	// Handle special case for "-" which means "ignore this field"
	if jsonTag == "-" {
		return field.Name
	}

	// Find the first comma to separate name from options
	if idx := strings.IndexByte(jsonTag, ','); idx != -1 {
		name := jsonTag[:idx]
		if name == "" || name == "-" {
			return field.Name
		}
		return name
	}

	// No comma found, return the entire tag
	return jsonTag
}
