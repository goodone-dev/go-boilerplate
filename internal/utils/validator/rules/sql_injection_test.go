package rules

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `validate:"sql_injection_safe"`
}

func TestSQLInjectionSafe(t *testing.T) {
	v := validator.New()
	_ = v.RegisterValidation("sql_injection_safe", SQLInjectionSafeValidation)

	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{
			name:     "Valid input",
			input:    "John Doe",
			hasError: false,
		},
		{
			name:     "Safe characters",
			input:    "O'Connor",
			hasError: false,
		},
		{
			name:     "SQL Injection - OR 1=1",
			input:    "John' OR '1'='1",
			hasError: true,
		},
		{
			name:     "SQL Injection - UNION SELECT",
			input:    "User UNION SELECT * FROM users",
			hasError: true,
		},
		{
			name:     "SQL Injection - Comment --",
			input:    "admin' --",
			hasError: true,
		},
		{
			name:     "SQL Injection - Comment /* */",
			input:    "/* Drop Table */",
			hasError: true,
		},
		{
			name:     "SQL Injection - Semicolon",
			input:    "bobby; DROP TABLE users",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := TestStruct{Name: tt.input}
			err := v.Struct(obj)
			if tt.hasError {
				assert.Error(t, err, "Expected validation errors but got none")
			} else {
				assert.NoError(t, err, "Expected no validation errors but got some")
			}
		})
	}
}
