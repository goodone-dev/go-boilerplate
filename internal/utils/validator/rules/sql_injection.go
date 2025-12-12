package rules

import (
	"regexp"

	universal "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func SQLInjectionSafeValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// Common SQL injection patterns:
	// - Union-based: UNION SELECT
	// - Error-based: ' OR '1'='1
	// - Comment-based: --, /* */
	// - Stacked queries: ;
	sqlInjectionPattern := `(?i)\b(UNION|SELECT|INSERT|UPDATE|DELETE|DROP|ALTER|CREATE|TRUNCATE|EXEC)\b|--|/\*|\*/|;|\bOR\b\s+.+?=.+`

	matched, _ := regexp.MatchString(sqlInjectionPattern, value)
	return !matched
}

func SQLInjectionSafeRegistration(ut universal.Translator) error {
	return ut.Add("sql_injection_safe", "{0} contains unsafe SQL characters", true)
}

func SQLInjectionSafeTranslation(ut universal.Translator, fe validator.FieldError) string {
	t, _ := ut.T("sql_injection_safe", fe.Field())
	return t
}
