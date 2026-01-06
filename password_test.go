package xvalidator

import (
	"strings"
	"testing"
)

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid strong password",
			password: "Test1234!",
			wantErr:  false,
		},
		{
			name:     "valid password with all character types",
			password: "Abcd1234!@#$",
			wantErr:  false,
		},
		{
			name:     "valid password with special characters",
			password: "MyP@ssw0rd!",
			wantErr:  false,
		},
		{
			name:     "valid password at minimum length",
			password: "Abc123!@",
			wantErr:  false,
		},
		{
			name:     "valid password near maximum length",
			password: "Abc123!" + strings.Repeat("x", 93),
			wantErr:  false,
		},
		{
			name:     "too short password",
			password: "Abc12!",
			wantErr:  true,
			errMsg:   "password must be at least 8 characters long",
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
			errMsg:   "password must be at least 8 characters long",
		},
		{
			name:     "too long password",
			password: strings.Repeat("A", 101) + "bc123!",
			wantErr:  true,
			errMsg:   "password must not exceed 100 characters",
		},
		{
			name:     "missing uppercase letter",
			password: "abcdefg123!",
			wantErr:  true,
			errMsg:   "uppercase letter",
		},
		{
			name:     "missing lowercase letter",
			password: "ABCDEFG123!",
			wantErr:  true,
			errMsg:   "lowercase letter",
		},
		{
			name:     "missing digit",
			password: "Abcdefgh!",
			wantErr:  true,
			errMsg:   "digit",
		},
		{
			name:     "missing special character",
			password: "Abcdefgh123",
			wantErr:  true,
			errMsg:   "special character",
		},
		{
			name:     "missing multiple requirements (no upper, no special)",
			password: "abcdefg123",
			wantErr:  true,
			errMsg:   "uppercase letter",
		},
		{
			name:     "missing all but lowercase",
			password: "abcdefgh",
			wantErr:  true,
			errMsg:   "uppercase letter",
		},
		{
			name:     "only digits",
			password: "12345678",
			wantErr:  true,
			errMsg:   "uppercase letter",
		},
		{
			name:     "only special characters",
			password: "!@#$%^&*",
			wantErr:  true,
			errMsg:   "uppercase letter",
		},
		{
			name:     "password with spaces and special char",
			password: "Test 123!",
			wantErr:  false,
		},
		{
			name:     "all special characters supported",
			password: "Test1234()_+-=[]{}|;:,.<>?",
			wantErr:  false,
		},
		{
			name:     "password with @ symbol",
			password: "Test@123",
			wantErr:  false,
		},
		{
			name:     "password with # symbol",
			password: "Test#123",
			wantErr:  false,
		},
		{
			name:     "password with $ symbol",
			password: "Test$123",
			wantErr:  false,
		},
		{
			name:     "password with % symbol",
			password: "Test%123",
			wantErr:  false,
		},
		{
			name:     "password with ^ symbol",
			password: "Test^123",
			wantErr:  false,
		},
		{
			name:     "password with & symbol",
			password: "Test&123",
			wantErr:  false,
		},
		{
			name:     "password with * symbol",
			password: "Test*123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePasswordStrength() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidatePasswordStrength() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePasswordStrength() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestValidatePasswordStrength_EdgeCases(t *testing.T) {
	t.Run("password exactly 8 characters", func(t *testing.T) {
		err := ValidatePasswordStrength("Test123!")
		if err != nil {
			t.Errorf("Expected no error for 8 character password, got: %v", err)
		}
	})

	t.Run("password exactly 100 characters", func(t *testing.T) {
		// Create a password with exactly 100 characters
		password := "Test123!" + strings.Repeat("a", 92)
		err := ValidatePasswordStrength(password)
		if err != nil {
			t.Errorf("Expected no error for 100 character password, got: %v", err)
		}
	})

	t.Run("password exactly 101 characters", func(t *testing.T) {
		// Create a password with exactly 101 characters
		password := "Test123!" + strings.Repeat("a", 93)
		err := ValidatePasswordStrength(password)
		if err == nil {
			t.Error("Expected error for 101 character password, got nil")
		}
		if !strings.Contains(err.Error(), "must not exceed 100 characters") {
			t.Errorf("Expected error about exceeding 100 characters, got: %v", err)
		}
	})

	t.Run("password with unicode characters", func(t *testing.T) {
		// Unicode characters should fail special character requirement
		password := "Test1234你好"
		err := ValidatePasswordStrength(password)
		if err == nil {
			t.Error("Expected error for password with unicode but no special chars")
		}
	})
}

func BenchmarkValidatePasswordStrength(b *testing.B) {
	passwords := []string{
		"Test1234!",
		"abcdefgh",
		"ABCDEFGH",
		"12345678",
		"Abc123!@",
		"VeryLongPasswordWithAllRequirements123!@#",
	}

	for _, password := range passwords {
		b.Run(password, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ValidatePasswordStrength(password)
			}
		})
	}
}
