package sanitizer

import (
	"github.com/go-sanitize/sanitize"
)

// Docs: https://github.com/go-sanitize/sanitize
var customSanitizer *sanitize.Sanitizer

func Sanitize(obj any) (err error) {
	if customSanitizer == nil {
		customSanitizer, err = sanitize.New()
		if err != nil {
			return err
		}
	}

	return customSanitizer.Sanitize(&obj)
}
