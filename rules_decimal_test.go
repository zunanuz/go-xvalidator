package xvalidator

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDecimalComparatorFunctions(t *testing.T) {
	// Test data
	val1 := decimal.NewFromFloat(100.50)
	val2 := decimal.NewFromFloat(50.25)
	val3 := decimal.NewFromFloat(100.50)

	tests := []struct {
		name     string
		fn       func(a, b *decimal.Decimal) bool
		a        *decimal.Decimal
		b        *decimal.Decimal
		expected bool
	}{
		{"decimalGreaterThan - true", decimalGreaterThan, &val1, &val2, true},
		{"decimalGreaterThan - false", decimalGreaterThan, &val2, &val1, false},
		{"decimalGreaterThanOrEqual - true (greater)", decimalGreaterThanOrEqual, &val1, &val2, true},
		{"decimalGreaterThanOrEqual - true (equal)", decimalGreaterThanOrEqual, &val1, &val3, true},
		{"decimalGreaterThanOrEqual - false", decimalGreaterThanOrEqual, &val2, &val1, false},
		{"decimalLessThan - true", decimalLessThan, &val2, &val1, true},
		{"decimalLessThan - false", decimalLessThan, &val1, &val2, false},
		{"decimalLessThanOrEqual - true (less)", decimalLessThanOrEqual, &val2, &val1, true},
		{"decimalLessThanOrEqual - true (equal)", decimalLessThanOrEqual, &val1, &val3, true},
		{"decimalLessThanOrEqual - false", decimalLessThanOrEqual, &val1, &val2, false},
		{"decimalEqual - true", decimalEqual, &val1, &val3, true},
		{"decimalEqual - false", decimalEqual, &val1, &val2, false},
		{"decimalNotEqual - true", decimalNotEqual, &val1, &val2, true},
		{"decimalNotEqual - false", decimalNotEqual, &val1, &val3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateDecimalOperation(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		param      string
		comparator func(a, b *decimal.Decimal) bool
		expected   bool
	}{
		{"valid greater than", "150.00", "100.00", decimalGreaterThan, true},
		{"invalid greater than", "50.00", "100.00", decimalGreaterThan, false},
		{"valid equal", "100.00", "100.00", decimalEqual, true},
		{"invalid equal", "100.01", "100.00", decimalEqual, false},
		{"invalid decimal value", "invalid", "100.00", decimalGreaterThan, false},
		{"invalid decimal param", "100.00", "invalid", decimalGreaterThan, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh validator for each test
			v := validator.New()

			// Register a temporary validator using validateDecimalOperation
			validatorFunc := validateDecimalOperation(tt.comparator)
			tempTag := "temp_decimal_test"

			v.RegisterValidation(tempTag, validatorFunc)

			// Test using v.Var
			err := v.Var(tt.value, tempTag+"="+tt.param)

			if tt.expected {
				assert.NoError(t, err, "Expected no error for value: %s, param: %s", tt.value, tt.param)
			} else {
				assert.Error(t, err, "Expected error for value: %s, param: %s", tt.value, tt.param)
			}
		})
	}
}

func TestDecimalTypeFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "valid decimal.Decimal",
			input:    decimal.NewFromFloat(123.45),
			expected: "123.45",
		},
		{
			name:     "string input",
			input:    "not a decimal",
			expected: nil,
		},
		{
			name:     "int input",
			input:    123,
			expected: nil,
		},
		{
			name:     "float input",
			input:    123.45,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a reflect.Value from the input
			fieldValue := reflect.ValueOf(tt.input)
			result := decimalTypeFunc(fieldValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseDecimalParams(t *testing.T) {
	tests := []struct {
		name              string
		param             string
		expectedPrecision int32
		expectedScale     int32
	}{
		{"empty param - default", "", DefaultPrecision, DefaultScale},
		{"scale only - 2", "2", DefaultPrecision, 2},
		{"scale only - 0", "0", DefaultPrecision, 0},
		{"scale only - 18", "18", DefaultPrecision, 18},
		{"precision:scale - 38:18", "38:18", 38, 18},
		{"precision:scale - 10:6", "10:6", 10, 6},
		{"precision:scale - 15:2", "15:2", 15, 2},
		{"precision:scale - 20:0", "20:0", 20, 0},
		{"invalid format - keep defaults", "invalid", DefaultPrecision, DefaultScale},
		{"invalid precision:scale - keep defaults", "abc:def", DefaultPrecision, DefaultScale},
		{"partial invalid - precision only", "10:abc", 10, DefaultScale},
		{"partial invalid - scale only", "abc:6", DefaultPrecision, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			precision, scale := parseDecimalParams(tt.param)
			assert.Equal(t, tt.expectedPrecision, precision)
			assert.Equal(t, tt.expectedScale, scale)
		})
	}
}

func TestValidateDecimalPrecisionScale(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		precision int32
		scale     int32
		expected  bool
	}{
		// Valid cases
		{"valid decimal - 123.45 with 10:2", "123.45", 10, 2, true},
		{"valid decimal - 123.456 with 10:3", "123.456", 10, 3, true},
		{"valid integer - 12345 with 10:2", "12345", 10, 2, true},
		{"valid zero - 0 with 10:2", "0", 10, 2, true},
		{"valid zero decimal - 0.00 with 10:2", "0.00", 10, 2, true},
		{"valid negative - -123.45 with 10:2", "-123.45", 10, 2, true},
		{"valid large number - 12345678.12 with 10:2", "12345678.12", 10, 2, true},
		{"valid no decimal part - 123 with 10:2", "123", 10, 2, true},

		// Scale validation
		{"valid scale 0 - 12345 with 10:0", "12345", 10, 0, true},
		{"invalid scale - 123.456 with 10:2", "123.456", 10, 2, false},
		{"invalid scale - 123.1 with 10:0", "123.1", 10, 0, false},

		// Precision validation (integer digits vs available space)
		{"valid precision edge - 12345678 with 10:2", "12345678", 10, 2, true},     // 8 integer digits, max allowed = 10-2 = 8
		{"invalid precision - 123456789 with 10:2", "123456789", 10, 2, false},     // 9 integer digits, max allowed = 10-2 = 8
		{"valid edge precision - 1234567890 with 10:0", "1234567890", 10, 0, true}, // 10 integer digits, max allowed = 10-0 = 10
		{"invalid precision - 12345678901 with 10:0", "12345678901", 10, 0, false}, // 11 integer digits, max allowed = 10-0 = 10

		// Default precision and scale (DefaultPrecision:DefaultScale)
		{"valid default - small number", "123.456789", DefaultPrecision, DefaultScale, true},
		{"valid default - large number", "12345678901234567890.123456789012345678", DefaultPrecision, DefaultScale, true},

		// Edge case for 5:2
		{"valid 5:2 - 123.45", "123.45", 5, 2, true},      // 3 integer digits, max allowed = 5-2 = 3
		{"invalid 5:2 - 1234.56", "1234.56", 5, 2, false}, // 4 integer digits, max allowed = 5-2 = 3
		{"valid 5:2 - 12345", "12345", 5, 2, false},       // 5 integer digits, max allowed = 5-2 = 3
		{"valid 5:2 edge - 123.4", "123.4", 5, 2, true},   // 3 integer digits, 1 decimal place
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := decimal.NewFromString(tt.value)
			assert.NoError(t, err)

			result := validateDecimalPrecisionScale(value, tt.precision, tt.scale)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateDecimal(t *testing.T) {
	// Setup validator
	v := validator.New()
	RegisterDecimalValidators(v)

	// Test struct for validation
	type TestStruct struct {
		Amount1 string `validate:"decimal"`       // default DefaultPrecision:DefaultScale
		Amount2 string `validate:"decimal=2"`     // DefaultPrecision:2
		Amount3 string `validate:"decimal=0"`     // DefaultPrecision:0 (integer only)
		Amount4 string `validate:"decimal=38:18"` // explicit 38:18
		Amount5 string `validate:"decimal=10:6"`  // custom 10:6
	}

	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
		field   string
	}{
		{
			name: "all valid - default format",
			input: TestStruct{
				Amount1: "123.456789012345678",
				Amount2: "123.45",
				Amount3: "12345",
				Amount4: "123.456789012345678",
				Amount5: "1234.567890",
			},
			wantErr: false,
		},
		{
			name: "invalid scale - Amount2",
			input: TestStruct{
				Amount1: "123.45",
				Amount2: "123.456", // should be scale=2
				Amount3: "12345",
				Amount4: "123.45",
				Amount5: "1234.56",
			},
			wantErr: true,
			field:   "Amount2",
		},
		{
			name: "invalid scale - Amount3 (integer only)",
			input: TestStruct{
				Amount1: "123.45",
				Amount2: "123.45",
				Amount3: "123.45", // should be integer (scale=0)
				Amount4: "123.45",
				Amount5: "1234.56",
			},
			wantErr: true,
			field:   "Amount3",
		},
		{
			name: "invalid precision - Amount5",
			input: TestStruct{
				Amount1: "123.45",
				Amount2: "123.45",
				Amount3: "12345",
				Amount4: "123.45",
				Amount5: "12345678901.567890", // too many digits (precision=10)
			},
			wantErr: true,
			field:   "Amount5",
		},
		{
			name: "invalid decimal format",
			input: TestStruct{
				Amount1: "not-a-number",
				Amount2: "123.45",
				Amount3: "12345",
				Amount4: "123.45",
				Amount5: "1234.56",
			},
			wantErr: true,
			field:   "Amount1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.field != "" {
					// Check if error is for the expected field
					validationErrors := err.(validator.ValidationErrors)
					found := false
					for _, fieldErr := range validationErrors {
						if fieldErr.Field() == tt.field {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected validation error for field %s", tt.field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDecimalWithDecimalType(t *testing.T) {
	// Setup validator
	v := validator.New()
	RegisterDecimalValidators(v)

	// Test struct with decimal.Decimal type
	type PaymentRequest struct {
		AmountTHB decimal.Decimal `json:"amount_thb" validate:"required,decimal=2"`
		Fee       decimal.Decimal `json:"fee" validate:"decimal=4"`
		Quantity  decimal.Decimal `json:"quantity" validate:"decimal=0"`
		Rate      decimal.Decimal `json:"rate" validate:"decimal=10:6"`
	}

	// Helper function to create decimal values
	mustDecimal := func(s string) decimal.Decimal {
		d, err := decimal.NewFromString(s)
		if err != nil {
			t.Fatalf("Failed to create decimal from %s: %v", s, err)
		}
		return d
	}

	tests := []struct {
		name    string
		input   PaymentRequest
		wantErr bool
		field   string
	}{
		{
			name: "all valid decimal.Decimal values",
			input: PaymentRequest{
				AmountTHB: mustDecimal("1234.56"),     // valid scale=2
				Fee:       mustDecimal("12.3456"),     // valid scale=4
				Quantity:  mustDecimal("100"),         // valid integer (scale=0)
				Rate:      mustDecimal("1234.567890"), // valid 10:6
			},
			wantErr: false,
		},
		{
			name: "invalid scale for AmountTHB",
			input: PaymentRequest{
				AmountTHB: mustDecimal("1234.567"), // invalid scale=3, should be 2
				Fee:       mustDecimal("12.34"),
				Quantity:  mustDecimal("100"),
				Rate:      mustDecimal("1234.56"),
			},
			wantErr: true,
			field:   "AmountTHB",
		},
		{
			name: "invalid scale for Quantity (should be integer)",
			input: PaymentRequest{
				AmountTHB: mustDecimal("1234.56"),
				Fee:       mustDecimal("12.34"),
				Quantity:  mustDecimal("100.5"), // invalid, should be integer
				Rate:      mustDecimal("1234.56"),
			},
			wantErr: true,
			field:   "Quantity",
		},
		{
			name: "invalid precision for Rate",
			input: PaymentRequest{
				AmountTHB: mustDecimal("1234.56"),
				Fee:       mustDecimal("12.34"),
				Quantity:  mustDecimal("100"),
				Rate:      mustDecimal("12345678901.567890"), // exceeds precision=10
			},
			wantErr: true,
			field:   "Rate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.field != "" {
					// Check if error is for the expected field
					validationErrors := err.(validator.ValidationErrors)
					found := false
					for _, fieldErr := range validationErrors {
						if fieldErr.Field() == tt.field {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected validation error for field %s", tt.field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseDecimalIfParam(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		expectedRule   string
		expectedField  string
		expectedExpect string
		expectError    bool
	}{
		// Valid cases
		{
			name:           "basic scale format - 2@Mode=mode1",
			param:          "2@Mode=mode1",
			expectedRule:   "2",
			expectedField:  "Mode",
			expectedExpect: "mode1",
			expectError:    false,
		},
		{
			name:           "precision:scale format - 38:19@Mode=mode2",
			param:          "38:19@Mode=mode2",
			expectedRule:   "38:19",
			expectedField:  "Mode",
			expectedExpect: "mode2",
			expectError:    false,
		},
		{
			name:           "integer only format - 0@Type=integer",
			param:          "0@Type=integer",
			expectedRule:   "0",
			expectedField:  "Type",
			expectedExpect: "integer",
			expectError:    false,
		},
		{
			name:           "empty rule format - @Status=active",
			param:          "@Status=active",
			expectedRule:   "",
			expectedField:  "Status",
			expectedExpect: "active",
			expectError:    false,
		},
		{
			name:           "complex field name - 10:6@PaymentType=credit_card",
			param:          "10:6@PaymentType=credit_card",
			expectedRule:   "10:6",
			expectedField:  "PaymentType",
			expectedExpect: "credit_card",
			expectError:    false,
		},

		// Invalid cases
		{
			name:        "missing @ separator - 2Mode=mode1",
			param:       "2Mode=mode1",
			expectError: true,
		},
		{
			name:        "missing = separator - 2@Modemode1",
			param:       "2@Modemode1",
			expectError: true,
		},
		{
			name:        "empty parameter",
			param:       "",
			expectError: true,
		},
		{
			name:        "only @ - @",
			param:       "@",
			expectError: true,
		},
		{
			name:        "only = - =",
			param:       "=",
			expectError: true,
		},
		{
			name:        "multiple @ - 2@Mode@=mode1",
			param:       "2@Mode@=mode1",
			expectError: true,
		},
		{
			name:        "multiple = - 2@Mode=mode=1",
			param:       "2@Mode=mode=1",
			expectError: true,
		},
		{
			name:           "missing field name - 2@=mode1",
			param:          "2@=mode1",
			expectError:    false, // This should parse but field will be empty
			expectedRule:   "2",
			expectedField:  "",
			expectedExpect: "mode1",
		},
		{
			name:           "missing expected value - 2@Mode=",
			param:          "2@Mode=",
			expectError:    false, // This should parse but expect will be empty
			expectedRule:   "2",
			expectedField:  "Mode",
			expectedExpect: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, field, expect, err := parseDecimalIfParam(tt.param)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRule, rule)
				assert.Equal(t, tt.expectedField, field)
				assert.Equal(t, tt.expectedExpect, expect)
			}
		})
	}
}

func TestValidateDecimalIf(t *testing.T) {
	// Setup validator
	v := validator.New()
	RegisterDecimalValidators(v)

	// Test struct for conditional decimal validation
	type PaymentRequest struct {
		Mode    string `json:"mode"`
		Type    string `json:"type"`
		Status  string `json:"status"`
		Amount1 string `json:"amount1" validate:"decimal_if=2@Mode=precise"`        // scale=2 if Mode="precise"
		Amount2 string `json:"amount2" validate:"decimal_if=0@Type=integer"`        // integer if Type="integer"
		Amount3 string `json:"amount3" validate:"decimal_if=10:6@Status=active"`    // 10:6 if Status="active"
		Amount4 string `json:"amount4" validate:"decimal_if=@Mode=default"`         // default if Mode="default"
		Fee     string `json:"fee" validate:"decimal_if=38:18@Mode=high_precision"` // high precision if Mode="high_precision"
	}

	tests := []struct {
		name    string
		input   PaymentRequest
		wantErr bool
		field   string
	}{
		// Valid cases - conditions met
		{
			name: "condition met - precise mode with valid scale",
			input: PaymentRequest{
				Mode:    "precise",
				Type:    "decimal",
				Status:  "inactive",
				Amount1: "123.45",     // valid scale=2 for Mode="precise"
				Amount2: "123.456",    // condition not met, skip validation
				Amount3: "123.456789", // condition not met, skip validation
				Amount4: "123.456789", // condition not met, skip validation
				Fee:     "123.456789", // condition not met, skip validation
			},
			wantErr: false,
		},
		{
			name: "condition met - integer type with valid integer",
			input: PaymentRequest{
				Mode:    "normal",
				Type:    "integer",
				Status:  "inactive",
				Amount1: "123.456",    // condition not met, skip validation
				Amount2: "12345",      // valid integer for Type="integer"
				Amount3: "123.456789", // condition not met, skip validation
				Amount4: "123.456789", // condition not met, skip validation
				Fee:     "123.456789", // condition not met, skip validation
			},
			wantErr: false,
		},
		{
			name: "condition met - active status with valid precision:scale",
			input: PaymentRequest{
				Mode:    "normal",
				Type:    "decimal",
				Status:  "active",
				Amount1: "123.456",     // condition not met, skip validation
				Amount2: "123.456",     // condition not met, skip validation
				Amount3: "1234.567890", // valid 10:6 for Status="active"
				Amount4: "123.456789",  // condition not met, skip validation
				Fee:     "123.456789",  // condition not met, skip validation
			},
			wantErr: false,
		},
		{
			name: "condition met - default mode with default precision:scale",
			input: PaymentRequest{
				Mode:    "default",
				Type:    "decimal",
				Status:  "inactive",
				Amount1: "123.456",          // condition not met, skip validation
				Amount2: "123.456",          // condition not met, skip validation
				Amount3: "123.456789",       // condition not met, skip validation
				Amount4: "123.456789012345", // valid default precision:scale for Mode="default"
				Fee:     "123.456789",       // condition not met, skip validation
			},
			wantErr: false,
		},
		{
			name: "condition met - high precision mode",
			input: PaymentRequest{
				Mode:    "high_precision",
				Type:    "decimal",
				Status:  "inactive",
				Amount1: "123.456",                                 // condition not met, skip validation
				Amount2: "123.456",                                 // condition not met, skip validation
				Amount3: "123.456789",                              // condition not met, skip validation
				Amount4: "123.456789",                              // condition not met, skip validation
				Fee:     "12345678901234567890.123456789012345678", // valid 38:18 for Mode="high_precision"
			},
			wantErr: false,
		},

		// Invalid cases - conditions met but validation fails
		{
			name: "condition met - precise mode with invalid scale",
			input: PaymentRequest{
				Mode:    "precise",
				Type:    "decimal",
				Status:  "inactive",
				Amount1: "123.456",    // invalid scale=3, should be 2 for Mode="precise"
				Amount2: "123.456",    // condition not met, skip validation
				Amount3: "123.456789", // condition not met, skip validation
				Amount4: "123.456789", // condition not met, skip validation
				Fee:     "123.456789", // condition not met, skip validation
			},
			wantErr: true,
			field:   "Amount1",
		},
		{
			name: "condition met - integer type with invalid decimal",
			input: PaymentRequest{
				Mode:    "normal",
				Type:    "integer",
				Status:  "inactive",
				Amount1: "123.456",    // condition not met, skip validation
				Amount2: "123.45",     // invalid decimal, should be integer for Type="integer"
				Amount3: "123.456789", // condition not met, skip validation
				Amount4: "123.456789", // condition not met, skip validation
				Fee:     "123.456789", // condition not met, skip validation
			},
			wantErr: true,
			field:   "Amount2",
		},
		{
			name: "condition met - active status with invalid precision",
			input: PaymentRequest{
				Mode:    "normal",
				Type:    "decimal",
				Status:  "active",
				Amount1: "123.456",            // condition not met, skip validation
				Amount2: "123.456",            // condition not met, skip validation
				Amount3: "12345678901.567890", // invalid precision, exceeds 10:6 for Status="active"
				Amount4: "123.456789",         // condition not met, skip validation
				Fee:     "123.456789",         // condition not met, skip validation
			},
			wantErr: true,
			field:   "Amount3",
		},

		// Edge cases
		{
			name: "no conditions met - all should pass",
			input: PaymentRequest{
				Mode:    "unknown",
				Type:    "unknown",
				Status:  "unknown",
				Amount1: "invalid_decimal", // condition not met, skip validation
				Amount2: "invalid_decimal", // condition not met, skip validation
				Amount3: "invalid_decimal", // condition not met, skip validation
				Amount4: "invalid_decimal", // condition not met, skip validation
				Fee:     "invalid_decimal", // condition not met, skip validation
			},
			wantErr: false,
		},
		{
			name: "invalid decimal format when condition met",
			input: PaymentRequest{
				Mode:    "precise",
				Type:    "decimal",
				Status:  "inactive",
				Amount1: "not-a-number", // invalid decimal format for Mode="precise"
				Amount2: "123.456",      // condition not met, skip validation
				Amount3: "123.456789",   // condition not met, skip validation
				Amount4: "123.456789",   // condition not met, skip validation
				Fee:     "123.456789",   // condition not met, skip validation
			},
			wantErr: true,
			field:   "Amount1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.field != "" {
					// Check if error is for the expected field
					validationErrors := err.(validator.ValidationErrors)
					found := false
					for _, fieldErr := range validationErrors {
						if fieldErr.Field() == tt.field {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected validation error for field %s", tt.field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
