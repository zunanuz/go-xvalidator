package xvalidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMobileE164 tests the mobile_e164 validation rule.
func TestMobileE164(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name        string
		phoneNumber string
		wantErr     bool
		description string
	}{
		{
			name:        "valid_thai_mobile_plus_prefix",
			phoneNumber: "+66812345678",
			wantErr:     false,
			description: "Thai mobile number with +66 prefix",
		},
		{
			name:        "valid_french_mobile_plus_prefix",
			phoneNumber: "+33612345678",
			wantErr:     false,
			description: "French mobile number with +33 prefix",
		},
		{
			name:        "valid_uk_mobile_plus_prefix",
			phoneNumber: "+447911123456",
			wantErr:     false,
			description: "UK mobile number with +44 prefix",
		},
		{
			name:        "invalid_format_no_plus",
			phoneNumber: "66812345678",
			wantErr:     true,
			description: "Thai number without + prefix (not E.164)",
		},
		{
			name:        "invalid_format_too_long",
			phoneNumber: "+661234567890123456",
			wantErr:     true,
			description: "Number too long (exceeds 15 digits)",
		},
		{
			name:        "invalid_format_too_short",
			phoneNumber: "+66123",
			wantErr:     true,
			description: "Number too short",
		},
		{
			name:        "invalid_format_letters",
			phoneNumber: "+66abc123456",
			wantErr:     true,
			description: "Number contains letters",
		},
		{
			name:        "invalid_empty_string",
			phoneNumber: "",
			wantErr:     true,
			description: "Empty phone number",
		},
		{
			name:        "invalid_plus_only",
			phoneNumber: "+",
			wantErr:     true,
			description: "Only plus sign",
		},
		{
			name:        "invalid_landline_number",
			phoneNumber: "+6621234567",
			wantErr:     true,
			description: "Thai landline number (should fail for mobile_e164)",
		},
		{
			name:        "invalid_special_chars",
			phoneNumber: "+66-812-345-678",
			wantErr:     true, // E164 regex rejects dashes (stricter validation)
			description: "Number with dashes (rejected by E164 regex)",
		},
		{
			name:        "invalid_spaces",
			phoneNumber: "+66 81 234 5678",
			wantErr:     true, // E164 regex rejects spaces (stricter validation)
			description: "Number with spaces (rejected by E164 regex)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test struct
			testStruct := struct {
				Phone string `validate:"mobile_e164"`
			}{
				Phone: tt.phoneNumber,
			}

			err := v.Validate(testStruct)

			if tt.wantErr {
				assert.Error(t, err, "Expected validation error for %s: %s", tt.name, tt.description)
			} else {
				assert.NoError(t, err, "Expected no validation error for %s: %s", tt.name, tt.description)
			}
		})
	}
}

// TestMobileE164WithE164Tag tests mobile_e164 rule used together with e164 tag.
func TestMobileE164WithE164Tag(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name        string
		phoneNumber string
		wantErr     bool
		description string
	}{
		{
			name:        "valid_mobile_with_e164_tag",
			phoneNumber: "+66812345678",
			wantErr:     false,
			description: "Valid mobile number in E.164 format",
		},
		{
			name:        "invalid_mobile_without_plus",
			phoneNumber: "66812345678",
			wantErr:     true,
			description: "Mobile number without + prefix (fails both e164 and mobile_e164)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test struct with both e164 and mobile_e164 tags
			testStruct := struct {
				Phone string `validate:"e164,mobile_e164"`
			}{
				Phone: tt.phoneNumber,
			}

			err := v.Validate(testStruct)

			if tt.wantErr {
				assert.Error(t, err, "Expected validation error for %s: %s", tt.name, tt.description)
			} else {
				assert.NoError(t, err, "Expected no validation error for %s: %s", tt.name, tt.description)
			}
		})
	}
}

// TestMobileE164EdgeCases tests edge cases for mobile_e164 validation.
func TestMobileE164EdgeCases(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name        string
		phoneNumber string
		wantErr     bool
		description string
	}{
		{
			name:        "valid_shortest_mobile",
			phoneNumber: "+66812345678",
			wantErr:     false,
			description: "Thai mobile number (valid example)",
		},
		{
			name:        "valid_australian_mobile",
			phoneNumber: "+61412345678",
			wantErr:     false,
			description: "Australian mobile number",
		},
		{
			name:        "invalid_country_code_000",
			phoneNumber: "+0001234567890",
			wantErr:     true,
			description: "Invalid country code 000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStruct := struct {
				Phone string `validate:"mobile_e164"`
			}{
				Phone: tt.phoneNumber,
			}

			err := v.Validate(testStruct)

			if tt.wantErr {
				assert.Error(t, err, "Expected validation error for %s: %s", tt.name, tt.description)
			} else {
				assert.NoError(t, err, "Expected no validation error for %s: %s", tt.name, tt.description)
			}
		})
	}
}

// TestMobileE164CountrySpecific tests country-specific mobile_e164 validation.
func TestMobileE164CountrySpecific(t *testing.T) {
	v, err := NewValidator()
	require.NoError(t, err)

	tests := []struct {
		name        string
		phoneNumber string
		country     string
		wantErr     bool
		description string
	}{
		{
			name:        "valid_thai_mobile_with_TH_param",
			phoneNumber: "+66812345678",
			country:     "TH",
			wantErr:     false,
			description: "Thai mobile number with TH parameter",
		},
		{
			name:        "valid_french_mobile_with_FR_param",
			phoneNumber: "+33612345678",
			country:     "FR",
			wantErr:     false,
			description: "French mobile number with FR parameter",
		},
		{
			name:        "valid_uk_mobile_with_GB_param",
			phoneNumber: "+447912345678",
			country:     "GB",
			wantErr:     false,
			description: "UK mobile number with GB parameter",
		},
		{
			name:        "invalid_thai_mobile_with_US_param",
			phoneNumber: "+66812345678",
			country:     "US",
			wantErr:     true,
			description: "Thai mobile number with wrong US parameter",
		},
		{
			name:        "invalid_us_mobile_with_TH_param",
			phoneNumber: "+19171234567",
			country:     "TH",
			wantErr:     true,
			description: "US mobile number with wrong TH parameter",
		},
		{
			name:        "invalid_uk_mobile_with_TH_param",
			phoneNumber: "+447912345678",
			country:     "TH",
			wantErr:     true,
			description: "UK mobile number with wrong TH parameter",
		},
		{
			name:        "invalid_landline_with_TH_param",
			phoneNumber: "+6621234567",
			country:     "TH",
			wantErr:     true,
			description: "Thai landline number with TH parameter (not mobile)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test struct dynamically based on country parameter
			var testStruct interface{}
			switch tt.country {
			case "TH":
				testStruct = struct {
					Phone string `validate:"mobile_e164=TH"`
				}{Phone: tt.phoneNumber}
			case "US":
				testStruct = struct {
					Phone string `validate:"mobile_e164=US"`
				}{Phone: tt.phoneNumber}
			case "GB":
				testStruct = struct {
					Phone string `validate:"mobile_e164=GB"`
				}{Phone: tt.phoneNumber}
			default:
				testStruct = struct {
					Phone string `validate:"mobile_e164"`
				}{Phone: tt.phoneNumber}
			}

			err := v.Validate(testStruct)

			if tt.wantErr {
				assert.Error(t, err, "Expected validation error for %s: %s", tt.name, tt.description)
			} else {
				assert.NoError(t, err, "Expected no validation error for %s: %s", tt.name, tt.description)
			}
		})
	}
}
