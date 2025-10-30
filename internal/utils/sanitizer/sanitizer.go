package sanitizer

import (
	"context"

	"github.com/go-sanitize/sanitize"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
)

// Docs: https://github.com/go-sanitize/sanitize
func NewSanitizer() *sanitize.Sanitizer {
	s, err := sanitize.New()
	if err != nil {
		logger.Fatal(context.Background(), err, "❌ Failed to create sanitizer")
		return nil
	}

	return s
}

var customSanitizer = NewSanitizer()

func Sanitize[S any](obj S) (err error) {
	return customSanitizer.Sanitize(&obj)
}
