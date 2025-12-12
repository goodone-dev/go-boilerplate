package masker

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type NestedStruct struct {
	Secret string `json:"secret"`
	Public string `json:"public"`
	Email  string `json:"email"`
}

type TestStruct struct {
	Password  string            `json:"password"`
	Token     string            `json:"token"`
	Data      NestedStruct      `json:"data"`
	DataPtr   *NestedStruct     `json:"data_ptr"`
	MapData   map[string]string `json:"map_data"`
	SliceData []NestedStruct    `json:"slice_data"`
	Age       int               `json:"age"`
	PIN       int               `json:"pin"` // Non-string sensitive
}

func TestMask(t *testing.T) {
	cases := []struct {
		name     string
		input    any
		expected func(t *testing.T, res any)
	}{
		{
			name: "Basic Struct",
			input: TestStruct{
				Password: "supersecretpassword",
				Token:    "abcdef123456",
				Age:      30,
				PIN:      1234,
			},
			expected: func(t *testing.T, res any) {
				m, ok := res.(map[string]any)
				assert.True(t, ok)
				assert.NotEqual(t, "supersecretpassword", m["password"])
				assert.NotEqual(t, "abcdef123456", m["token"])
				assert.Equal(t, 30, m["age"])
				assert.Equal(t, "******", m["pin"]) // Non-string sensitive should be masked
			},
		},
		{
			name: "Nested Struct",
			input: NestedStruct{
				Secret: "internal_secret",
				Public: "public_info",
				Email:  "user@example.com",
			},
			expected: func(t *testing.T, res any) {
				m, ok := res.(map[string]any)
				assert.True(t, ok)
				assert.NotEqual(t, "internal_secret", m["secret"])
				assert.Equal(t, "public_info", m["public"])
				// Check email masking format (partial mask)
				assert.NotEqual(t, "user@example.com", m["email"])
				assert.Contains(t, m["email"], "*")
			},
		},
		{
			name: "Pointer",
			input: &NestedStruct{
				Secret: "ptr_secret",
			},
			expected: func(t *testing.T, res any) {
				m, ok := res.(map[string]any)
				assert.True(t, ok)
				assert.NotEqual(t, "ptr_secret", m["secret"])
			},
		},
		{
			name: "Map",
			input: map[string]any{
				"password": "map_password",
				"normal":   "normal_value",
			},
			expected: func(t *testing.T, res any) {
				m, ok := res.(map[string]any)
				assert.True(t, ok)
				assert.NotEqual(t, "map_password", m["password"])
				assert.Equal(t, "normal_value", m["normal"])
			},
		},
		{
			name:  "Nil",
			input: nil,
			expected: func(t *testing.T, res any) {
				assert.Nil(t, res)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Mask(tc.input)

			// Debug output
			b, _ := json.MarshalIndent(result, "", "  ")
			t.Logf("Result: %s", string(b))

			tc.expected(t, result)
		})
	}
}
