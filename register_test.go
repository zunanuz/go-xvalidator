package xvalidator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRegisterDecimalValidators(t *testing.T) {
	v := validator.New()

	// Register decimal validators
	RegisterDecimalValidators(v)

	tests := []struct {
		name    string
		value   string
		tag     string
		wantErr bool
	}{
		{
			name:    "dgt - greater than valid",
			value:   "150.00",
			tag:     "dgt=100.00",
			wantErr: false,
		},
		{
			name:    "dgt - greater than invalid",
			value:   "50.00",
			tag:     "dgt=100.00",
			wantErr: true,
		},
		{
			name:    "dgte - greater than or equal valid",
			value:   "100.00",
			tag:     "dgte=100.00",
			wantErr: false,
		},
		{
			name:    "dgte - greater than or equal invalid",
			value:   "99.99",
			tag:     "dgte=100.00",
			wantErr: true,
		},
		{
			name:    "dlt - less than valid",
			value:   "50.00",
			tag:     "dlt=100.00",
			wantErr: false,
		},
		{
			name:    "dlt - less than invalid",
			value:   "150.00",
			tag:     "dlt=100.00",
			wantErr: true,
		},
		{
			name:    "dlte - less than or equal valid",
			value:   "100.00",
			tag:     "dlte=100.00",
			wantErr: false,
		},
		{
			name:    "dlte - less than or equal invalid",
			value:   "100.01",
			tag:     "dlte=100.00",
			wantErr: true,
		},
		{
			name:    "deq - equal valid",
			value:   "100.00",
			tag:     "deq=100.00",
			wantErr: false,
		},
		{
			name:    "deq - equal invalid",
			value:   "100.01",
			tag:     "deq=100.00",
			wantErr: true,
		},
		{
			name:    "dneq - not equal valid",
			value:   "100.01",
			tag:     "dneq=100.00",
			wantErr: false,
		},
		{
			name:    "dneq - not equal invalid",
			value:   "100.00",
			tag:     "dneq=100.00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.value, tt.tag)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegisterDecimalValidators_DecimalType(t *testing.T) {
	v := validator.New()
	RegisterDecimalValidators(v)

	type TestStruct struct {
		Amount decimal.Decimal `validate:"required"`
	}

	// Test that decimal.Decimal type is properly registered
	data := TestStruct{
		Amount: decimal.NewFromFloat(100.50),
	}

	err := v.Struct(data)
	assert.NoError(t, err)
}

func TestRegisterDecimalValidators_InvalidDecimalString(t *testing.T) {
	v := validator.New()
	RegisterDecimalValidators(v)

	// Test with invalid decimal string
	err := v.Var("invalid_decimal", "dgt=100.00")
	assert.Error(t, err)
}

func TestRegisterDecimalValidators_InvalidParameterString(t *testing.T) {
	v := validator.New()
	RegisterDecimalValidators(v)

	// Test with invalid parameter string
	err := v.Var("100.00", "dgt=invalid_param")
	assert.Error(t, err)
}

func TestRegisterValidators_Integration(t *testing.T) {
	v := validator.New()

	// Register all validators
	RegisterDecimalValidators(v)

	type TestStruct struct {
		Amount string `validate:"dgt=100.00"`
	}

	// Test valid data
	validData := TestStruct{
		Amount: "150.00",
	}

	err := v.Struct(validData)
	assert.NoError(t, err)

	// Test invalid data
	invalidData := TestStruct{
		Amount: "50.00", // too small
	}

	err = v.Struct(invalidData)
	assert.Error(t, err)
}
