package validator

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// use a single instance , it caches struct info
var (
	uni        *ut.UniversalTranslator
	validate   *validator.Validate
	Translator *ut.Translator
)

func init() {

	// NOTE: ommitting allot of error checking for brevity

	en := en.New()
	uni = ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")
	Translator = &trans
	validate = validator.New()
	err := en_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic("unable to register translation")
	}

	err = validate.RegisterTranslation("required", *Translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	if err != nil {
		panic("unable to register translation")
	}

	err = validate.RegisterTranslation("alpha", *Translator, func(ut ut.Translator) error {
		return ut.Add("alpha", "{0} should be alphanumeric only", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alpha", fe.Field())

		return t
	})
	if err != nil {
		panic("unable to register translation: alpha")
	}

}

func Validate(s interface{}) error {

	var errors []string

	err := validate.Struct(s)
	if err != nil {

		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			errors = append(errors, e.Translate(*Translator))
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf(errors[0])

}
