package xvalidator

import (
	"testing"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_setupTranslator(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful translator setup",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()

			// Register custom validators first
			RegisterDecimalValidators(v)
			RegisterURLValidators(v)
			RegisterPhoneValidators(v)

			trans, err := setupTranslator(v)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, trans)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, trans)
			}
		})
	}
}

func Test_formatTranslatedErrors(t *testing.T) {
	// Setup validator with translator
	v := validator.New()
	RegisterDecimalValidators(v)
	RegisterURLValidators(v)
	RegisterPhoneValidators(v)

	trans, err := setupTranslator(v)
	require.NoError(t, err)
	require.NotNil(t, trans)

	// Test struct for validation errors
	type TestStruct struct {
		Email    string `validate:"required,email" json:"email"`
		Age      int    `validate:"required,min=18" json:"age"`
		Username string `validate:"required,min=3" json:"username"`
	}

	tests := []struct {
		name           string
		input          TestStruct
		wantErr        bool
		expectedErrors []string
	}{
		{
			name: "multiple validation errors",
			input: TestStruct{
				Email:    "invalid-email",
				Age:      15,
				Username: "ab",
			},
			wantErr: true,
			expectedErrors: []string{
				"Email must be a valid email address",
				"Age must be 18 or greater",
				"Username must be at least 3 characters in length",
			},
		},
		{
			name: "single validation error",
			input: TestStruct{
				Email:    "test@example.com",
				Age:      25,
				Username: "ab",
			},
			wantErr: true,
			expectedErrors: []string{
				"Username must be at least 3 characters in length",
			},
		},
		{
			name: "no validation errors",
			input: TestStruct{
				Email:    "test@example.com",
				Age:      25,
				Username: "validuser",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)

			if tt.wantErr {
				require.Error(t, err)

				validationErrors, ok := err.(validator.ValidationErrors)
				require.True(t, ok)

				translatedErr := formatTranslatedErrors(validationErrors, trans)
				assert.Error(t, translatedErr)

				errorMsg := translatedErr.Error()
				for _, expectedError := range tt.expectedErrors {
					assert.Contains(t, errorMsg, expectedError)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_registerDecimalTranslation(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "register decimal translation on fresh validator",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			RegisterDecimalValidators(v)

			// Create fresh translator without setupTranslator which already registers custom translations
			en := en.New()
			uni := ut.New(en, en)
			trans, _ := uni.GetTranslator("en")

			// Register default English translations only
			err := en_trans.RegisterDefaultTranslations(v, trans)
			require.NoError(t, err)

			err = registerDecimalTranslation(v, trans)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_registerDecimalIfTranslation(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "register decimal_if translation on fresh validator",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			RegisterDecimalValidators(v)

			// Create fresh translator without setupTranslator which already registers custom translations
			en := en.New()
			uni := ut.New(en, en)
			trans, _ := uni.GetTranslator("en")

			// Register default English translations only
			err := en_trans.RegisterDefaultTranslations(v, trans)
			require.NoError(t, err)

			err = registerDecimalIfTranslation(v, trans)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_registerCustomTranslations(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "register custom translations on fresh validator",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()

			// Register custom validators first
			RegisterDecimalValidators(v)
			RegisterURLValidators(v)
			RegisterPhoneValidators(v)

			// Create fresh translator without setupTranslator which already registers custom translations
			en := en.New()
			uni := ut.New(en, en)
			trans, _ := uni.GetTranslator("en")

			// Register default English translations only
			err := en_trans.RegisterDefaultTranslations(v, trans)
			require.NoError(t, err)

			err = registerCustomTranslations(v, trans)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDecimalTranslationMessages(t *testing.T) {
	// Setup validator with translator
	validator, err := NewValidator()
	require.NoError(t, err)
	require.NotNil(t, validator)

	type TestStruct struct {
		Amount       string `validate:"decimal=10:2" json:"amount"`
		IntegerValue string `validate:"decimal=0" json:"integer_value"`
		DefaultValue string `validate:"decimal" json:"default_value"`
	}

	tests := []struct {
		name           string
		input          TestStruct
		wantErr        bool
		expectedErrors []string
	}{
		{
			name: "invalid decimal with precision scale",
			input: TestStruct{
				Amount: "invalid",
			},
			wantErr: true,
			expectedErrors: []string{
				"amount must be a decimal with precision ≤ 10 and scale ≤ 2",
			},
		},
		{
			name: "invalid integer format",
			input: TestStruct{
				IntegerValue: "123.45",
			},
			wantErr: true,
			expectedErrors: []string{
				"integer_value must be an integer format (no decimal places)",
			},
		},
		{
			name: "invalid default decimal",
			input: TestStruct{
				DefaultValue: "invalid",
			},
			wantErr: true,
			expectedErrors: []string{
				"default_value must be a decimal with precision ≤ 38 and scale ≤ 18",
			},
		},
		{
			name: "valid decimal values",
			input: TestStruct{
				Amount:       "123.45",
				IntegerValue: "123",
				DefaultValue: "123.456789",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.StructTranslated(tt.input)

			if tt.wantErr {
				require.Error(t, err)

				errorMsg := err.Error()
				for _, expectedError := range tt.expectedErrors {
					assert.Contains(t, errorMsg, expectedError)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDecimalIfTranslationMessages(t *testing.T) {
	// Setup validator with translator
	validator, err := NewValidator()
	require.NoError(t, err)
	require.NotNil(t, validator)

	type TestStruct struct {
		PaymentType string `json:"payment_type"`
		Amount      string `validate:"decimal_if=10:2@PaymentType=credit" json:"amount"`
		Quantity    string `validate:"decimal_if=0@PaymentType=bulk" json:"quantity"`
	}

	tests := []struct {
		name           string
		input          TestStruct
		wantErr        bool
		expectedErrors []string
	}{
		{
			name: "invalid decimal_if with precision scale",
			input: TestStruct{
				PaymentType: "credit",
				Amount:      "invalid",
			},
			wantErr: true,
			expectedErrors: []string{
				"amount must be a decimal with precision ≤ 10 and scale ≤ 2 when PaymentType equals 'credit'",
			},
		},
		{
			name: "invalid decimal_if integer format",
			input: TestStruct{
				PaymentType: "bulk",
				Quantity:    "123.45",
			},
			wantErr: true,
			expectedErrors: []string{
				"quantity must be an integer format (no decimal places) when PaymentType equals 'bulk'",
			},
		},
		{
			name: "valid conditional decimal values",
			input: TestStruct{
				PaymentType: "credit",
				Amount:      "123.45",
			},
			wantErr: false,
		},
		{
			name: "condition not met - no validation",
			input: TestStruct{
				PaymentType: "debit",
				Amount:      "invalid",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.StructTranslated(tt.input)

			if tt.wantErr {
				require.Error(t, err)

				errorMsg := err.Error()
				for _, expectedError := range tt.expectedErrors {
					assert.Contains(t, errorMsg, expectedError)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCustomValidatorTranslationMessages(t *testing.T) {
	// Setup validator with translator
	validator, err := NewValidator()
	require.NoError(t, err)
	require.NotNil(t, validator)

	type TestStruct struct {
		Price       string `validate:"dgt=100" json:"price"`
		MinAge      string `validate:"dgte=18" json:"min_age"`
		MaxAge      string `validate:"dlt=65" json:"max_age"`
		Score       string `validate:"dlte=100" json:"score"`
		ExactValue  string `validate:"deq=50" json:"exact_value"`
		NotValue    string `validate:"dneq=0" json:"not_value"`
		WebsiteURL  string `validate:"https_url" json:"website_url"`
		PhoneNumber string `validate:"mobile_e164" json:"phone_number"`
	}

	tests := []struct {
		name           string
		input          TestStruct
		wantErr        bool
		expectedErrors []string
	}{
		{
			name: "decimal greater than validation",
			input: TestStruct{
				Price: "50",
			},
			wantErr: true,
			expectedErrors: []string{
				"price must be greater than 100",
			},
		},
		{
			name: "decimal greater than or equal validation",
			input: TestStruct{
				MinAge: "15",
			},
			wantErr: true,
			expectedErrors: []string{
				"min_age must be greater than or equal to 18",
			},
		},
		{
			name: "decimal less than validation",
			input: TestStruct{
				MaxAge: "70",
			},
			wantErr: true,
			expectedErrors: []string{
				"max_age must be less than 65",
			},
		},
		{
			name: "decimal less than or equal validation",
			input: TestStruct{
				Score: "150",
			},
			wantErr: true,
			expectedErrors: []string{
				"score must be less than or equal to 100",
			},
		},
		{
			name: "decimal equal validation",
			input: TestStruct{
				ExactValue: "45",
			},
			wantErr: true,
			expectedErrors: []string{
				"exact_value must be equal to 50",
			},
		},
		{
			name: "decimal not equal validation",
			input: TestStruct{
				NotValue: "0",
			},
			wantErr: true,
			expectedErrors: []string{
				"not_value must not be equal to 0",
			},
		},
		{
			name: "https url validation",
			input: TestStruct{
				WebsiteURL: "http://example.com",
			},
			wantErr: true,
			expectedErrors: []string{
				"website_url must be a valid HTTPS URL",
			},
		},
		{
			name: "mobile e164 validation",
			input: TestStruct{
				PhoneNumber: "0812345678",
			},
			wantErr: true,
			expectedErrors: []string{
				"phone_number must be a valid mobile number in E.164 format (e.g., +66812345678)",
			},
		},
		{
			name: "valid custom validator values",
			input: TestStruct{
				Price:       "150",
				MinAge:      "25",
				MaxAge:      "55",
				Score:       "85",
				ExactValue:  "50",
				NotValue:    "100",
				WebsiteURL:  "https://example.com",
				PhoneNumber: "+66812345678",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.StructTranslated(tt.input)

			if tt.wantErr {
				require.Error(t, err)

				errorMsg := err.Error()
				for _, expectedError := range tt.expectedErrors {
					assert.Contains(t, errorMsg, expectedError)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVarTranslatedWithCustomValidators(t *testing.T) {
	// Setup validator with translator
	validator, err := NewValidator()
	require.NoError(t, err)
	require.NotNil(t, validator)

	tests := []struct {
		name          string
		value         string
		tag           string
		wantErr       bool
		expectedError string
	}{
		{
			name:          "decimal validation with var",
			value:         "invalid",
			tag:           "decimal=10:2",
			wantErr:       true,
			expectedError: " must be a decimal with precision ≤ 10 and scale ≤ 2",
		},
		{
			name:          "https url validation with var",
			value:         "http://example.com",
			tag:           "https_url",
			wantErr:       true,
			expectedError: " must be a valid HTTPS URL",
		},
		{
			name:          "mobile e164 validation with var",
			value:         "0812345678",
			tag:           "mobile_e164",
			wantErr:       true,
			expectedError: " must be a valid mobile number in E.164 format (e.g., +66812345678)",
		},
		{
			name:    "valid decimal with var",
			value:   "123.45",
			tag:     "decimal=10:2",
			wantErr: false,
		},
		{
			name:    "valid https url with var",
			value:   "https://example.com",
			tag:     "https_url",
			wantErr: false,
		},
		{
			name:    "valid mobile number with var",
			value:   "+66812345678",
			tag:     "mobile_e164",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.VarTranslated(tt.value, tt.tag)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
