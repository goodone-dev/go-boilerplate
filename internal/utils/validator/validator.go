package validator

import (
	"fmt"

	"github.com/go-playground/locales/en"
	universal "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

type CustomValidator struct {
	validator  *validator.Validate
	translator universal.Translator
}

var customValidator *CustomValidator

func NewValidator() error {
	en := en.New()
	un := universal.New(en, en)

	vl := validator.New()
	tr, ok := un.GetTranslator("en")
	if !ok {
		return fmt.Errorf("english translator not found in universal translator")
	}

	err := translations.RegisterDefaultTranslations(vl, tr)
	if err != nil {
		return err
	}

	customValidator = &CustomValidator{
		validator:  vl,
		translator: tr,
	}

	return nil
}

func Validate(obj any) []string {
	if customValidator == nil {
		err := NewValidator()
		if err != nil {
			return []string{err.Error()}
		}
	}

	if err := customValidator.validator.Struct(obj); err != nil {
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Translate(customValidator.translator))
		}

		return errors
	}

	return nil
}
