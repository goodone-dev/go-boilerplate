package validator

import (
	"context"

	"github.com/go-playground/locales/en"
	universal "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/utils/validator/rules"
)

type CustomValidator struct {
	validator  *validator.Validate
	translator universal.Translator
}

func NewValidator() *CustomValidator {
	en := en.New()
	un := universal.New(en, en)

	vl := validator.New()
	tr, ok := un.GetTranslator("en")
	if !ok {
		logger.Fatal(context.Background(), nil, "❌ Failed to get translator")
		return nil
	}

	err := translations.RegisterDefaultTranslations(vl, tr)
	if err != nil {
		logger.Fatal(context.Background(), err, "❌ Failed to register default translations")
		return nil
	}

	// Register custom validators
	_ = vl.RegisterValidation("sql_injection_safe", rules.SQLInjectionSafeValidation)

	// Register custom translations
	_ = vl.RegisterTranslation("sql_injection_safe", tr, rules.SQLInjectionSafeRegistration, rules.SQLInjectionSafeTranslation)

	return &CustomValidator{
		validator:  vl,
		translator: tr,
	}
}

var customValidator = NewValidator()

func Validate(obj any) []string {
	if err := customValidator.validator.Struct(obj); err != nil {
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Translate(customValidator.translator))
		}

		return errors
	}

	return nil
}
