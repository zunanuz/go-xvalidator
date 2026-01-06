package xvalidator

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test structs for validation
type TestUser struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,min=18"`
	Website  string `json:"website,omitempty" validate:"omitempty,url"`
	Optional string `json:"optional,omitempty"`
}

type TestUserWithJSONTags struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	IgnoredID string `json:"-"`
	NoTag     string `validate:"required"`
}

type InvalidStruct struct {
	Field string `validate:"invalid_tag"`
}

func TestNewValidator(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful validator creation",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewValidator()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, v)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.NotNil(t, v.validate)
				assert.NotNil(t, v.translator)
			}
		})
	}
}

func TestValidator_Validate(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name: "valid struct",
			input: TestUser{
				Name:    "John Doe",
				Email:   "john@example.com",
				Age:     25,
				Website: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid struct - missing required fields",
			input: TestUser{
				Name: "", // required field empty
				Age:  17, // below minimum
			},
			wantErr: true,
		},
		{
			name: "invalid struct - invalid email",
			input: TestUser{
				Name:  "John",
				Email: "invalid-email",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "valid struct with optional fields",
			input: TestUser{
				Name:     "Jane",
				Email:    "jane@example.com",
				Age:      30,
				Optional: "some value",
			},
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "non-struct input",
			input:   "string",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_Struct(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name: "valid struct",
			input: TestUser{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "invalid struct",
			input: TestUser{
				Name: "J", // too short
				Age:  17,  // below minimum
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_Var(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name    string
		field   any
		tag     string
		wantErr bool
	}{
		{
			name:    "valid email",
			field:   "test@example.com",
			tag:     "email",
			wantErr: false,
		},
		{
			name:    "invalid email",
			field:   "invalid-email",
			tag:     "email",
			wantErr: true,
		},
		{
			name:    "valid required field",
			field:   "some value",
			tag:     "required",
			wantErr: false,
		},
		{
			name:    "empty required field",
			field:   "",
			tag:     "required",
			wantErr: true,
		},
		{
			name:    "valid min length",
			field:   "hello",
			tag:     "min=3",
			wantErr: false,
		},
		{
			name:    "invalid min length",
			field:   "hi",
			tag:     "min=3",
			wantErr: true,
		},
		{
			name:    "valid number range",
			field:   25,
			tag:     "min=18,max=65",
			wantErr: false,
		},
		{
			name:    "invalid number range",
			field:   17,
			tag:     "min=18,max=65",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.field, tt.tag)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_StructTranslated(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name           string
		input          any
		wantErr        bool
		checkErrMsg    bool
		expectedErrMsg string
	}{
		{
			name: "valid struct",
			input: TestUser{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "invalid struct with translated errors",
			input: TestUser{
				Name:  "", // required field empty
				Email: "invalid-email",
				Age:   17, // below minimum
			},
			wantErr:     true,
			checkErrMsg: true,
		},
		{
			name: "struct with multiple validation errors",
			input: TestUser{
				Name:    "J", // too short
				Email:   "",  // required and invalid
				Age:     17,  // below minimum
				Website: "not-a-url",
			},
			wantErr:     true,
			checkErrMsg: true,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.StructTranslated(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.checkErrMsg {
					// Check that error message is user-friendly (not raw validator error)
					assert.NotContains(t, err.Error(), "ValidationErrors")
					assert.NotEmpty(t, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_VarTranslated(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name        string
		field       any
		tag         string
		wantErr     bool
		checkErrMsg bool
	}{
		{
			name:    "valid email",
			field:   "test@example.com",
			tag:     "email",
			wantErr: false,
		},
		{
			name:        "invalid email with translated error",
			field:       "invalid-email",
			tag:         "email",
			wantErr:     true,
			checkErrMsg: true,
		},
		{
			name:        "empty required field with translated error",
			field:       "",
			tag:         "required",
			wantErr:     true,
			checkErrMsg: true,
		},
		{
			name:        "invalid min length with translated error",
			field:       "hi",
			tag:         "min=3",
			wantErr:     true,
			checkErrMsg: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.VarTranslated(tt.field, tt.tag)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.checkErrMsg {
					// Check that error message is user-friendly
					assert.NotEmpty(t, err.Error())
					assert.NotContains(t, err.Error(), "ValidationErrors")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetJSONTagName(t *testing.T) {
	tests := []struct {
		name     string
		field    reflect.StructField
		expected string
	}{
		{
			name: "field with simple json tag",
			field: reflect.StructField{
				Name: "TestField",
				Tag:  `json:"test_field"`,
			},
			expected: "test_field",
		},
		{
			name: "field with json tag and omitempty",
			field: reflect.StructField{
				Name: "TestField",
				Tag:  `json:"test_field,omitempty"`,
			},
			expected: "test_field",
		},
		{
			name: "field with json tag and multiple options",
			field: reflect.StructField{
				Name: "TestField",
				Tag:  `json:"test_field,omitempty,string"`,
			},
			expected: "test_field",
		},
		{
			name: "field with json ignore tag",
			field: reflect.StructField{
				Name: "TestField",
				Tag:  `json:"-"`,
			},
			expected: "TestField",
		},
		{
			name: "field with empty json tag",
			field: reflect.StructField{
				Name: "TestField",
				Tag:  `json:""`,
			},
			expected: "TestField",
		},
		{
			name: "field with no json tag",
			field: reflect.StructField{
				Name: "TestField",
				Tag:  `validate:"required"`,
			},
			expected: "TestField",
		},
		{
			name: "field with empty json name and options",
			field: reflect.StructField{
				Name: "TestField",
				Tag:  `json:",omitempty"`,
			},
			expected: "TestField",
		},
		{
			name: "field with dash and options",
			field: reflect.StructField{
				Name: "TestField",
				Tag:  `json:"-,omitempty"`,
			},
			expected: "TestField",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getJSONTagName(tt.field)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidator_ComparisonBetweenRawAndTranslated(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	invalidUser := TestUser{
		Name:  "",
		Email: "invalid-email",
		Age:   17,
	}

	// Test raw validation
	rawErr := v.Validate(invalidUser)
	require.Error(t, rawErr)

	// Test translated validation
	translatedErr := v.StructTranslated(invalidUser)
	require.Error(t, translatedErr)

	// Errors should be different types
	_, isValidationErrors := rawErr.(validator.ValidationErrors)
	assert.True(t, isValidationErrors, "Raw error should be ValidationErrors type")

	_, isValidationErrorsTranslated := translatedErr.(validator.ValidationErrors)
	assert.False(t, isValidationErrorsTranslated, "Translated error should NOT be ValidationErrors type")

	// Both should contain error information
	assert.NotEmpty(t, rawErr.Error())
	assert.NotEmpty(t, translatedErr.Error())
}

func TestValidator_StructVsValidateConsistency(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	testCases := []any{
		TestUser{Name: "John", Email: "john@example.com", Age: 25},
		TestUser{Name: "", Email: "invalid", Age: 17},
		TestUserWithJSONTags{FirstName: "John", LastName: "Doe", NoTag: "value"},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			validateErr := v.Validate(testCase)
			structErr := v.Struct(testCase)

			// Both methods should return the same result
			if validateErr == nil {
				assert.NoError(t, structErr)
			} else {
				assert.Error(t, structErr)
			}
		})
	}
}

func TestValidator_ErrorTypes(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	t.Run("validation errors are properly typed", func(t *testing.T) {
		invalidUser := TestUser{Name: "", Email: "invalid"}

		// Raw validation should return ValidationErrors
		rawErr := v.Validate(invalidUser)
		require.Error(t, rawErr)
		_, ok := rawErr.(validator.ValidationErrors)
		assert.True(t, ok, "Raw validation should return ValidationErrors")

		// Translated validation should return formatted error
		translatedErr := v.StructTranslated(invalidUser)
		require.Error(t, translatedErr)
		_, ok = translatedErr.(validator.ValidationErrors)
		assert.False(t, ok, "Translated validation should not return ValidationErrors")
	})

	t.Run("non-validation errors are passed through", func(t *testing.T) {
		// This would cause a different type of error (not ValidationErrors)
		err := v.Validate("not a struct")
		require.Error(t, err)

		translatedErr := v.StructTranslated("not a struct")
		require.Error(t, translatedErr)

		// Both should be the same error since it's not a ValidationErrors
		assert.Equal(t, err.Error(), translatedErr.Error())
	})
}
