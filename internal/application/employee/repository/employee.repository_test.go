package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmployeeRepository(t *testing.T) {
	// Test that the constructor doesn't panic with nil
	// In real usage, baseRepo would be properly initialized
	assert.NotPanics(t, func() {
		NewEmployeeRepository(nil)
	})
}
