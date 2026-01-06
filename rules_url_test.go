package xvalidator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateHttpsScheme(t *testing.T) {
	v := validator.New()
	RegisterURLValidators(v)

	type testStruct struct {
		URL string `validate:"https_url"`
	}

	tests := []struct {
		name    string
		input   testStruct
		wantErr bool
	}{
		{
			name:    "valid https url",
			input:   testStruct{URL: "https://www.example.com"},
			wantErr: false,
		},
		{
			name:    "valid https url with path",
			input:   testStruct{URL: "https://api.example.com/v1/users"},
			wantErr: false,
		},
		{
			name:    "valid https url with query parameters",
			input:   testStruct{URL: "https://example.com/search?q=test&limit=10"},
			wantErr: false,
		},
		{
			name:    "invalid http url",
			input:   testStruct{URL: "http://www.example.com"},
			wantErr: true,
		},
		{
			name:    "invalid ftp url",
			input:   testStruct{URL: "ftp://files.example.com"},
			wantErr: true,
		},
		{
			name:    "invalid url without protocol",
			input:   testStruct{URL: "www.example.com"},
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   testStruct{URL: ""},
			wantErr: true,
		},
		{
			name:    "invalid url format",
			input:   testStruct{URL: "not-a-url"},
			wantErr: true,
		},
		{
			name:    "valid https localhost",
			input:   testStruct{URL: "https://localhost:8080"},
			wantErr: false,
		},
		{
			name:    "invalid http localhost",
			input:   testStruct{URL: "http://localhost:8080"},
			wantErr: true,
		},
		{
			name:    "https url without host",
			input:   testStruct{URL: "https://"},
			wantErr: true,
		},
		{
			name:    "https scheme only",
			input:   testStruct{URL: "https:"},
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
