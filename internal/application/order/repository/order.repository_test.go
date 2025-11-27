package repository

import (
	"os"
	"testing"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.Disabled()
	code := m.Run()

	os.Exit(code)
}

func TestNewOrderRepository(t *testing.T) {
	// Test that the constructor doesn't panic with nil
	// In real usage, baseRepo would be properly initialized
	assert.NotPanics(t, func() {
		NewOrderRepository(nil)
	})
}
