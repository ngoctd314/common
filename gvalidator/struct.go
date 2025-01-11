package gvalidator

import (
	"errors"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// use a single instance, it caches struct info
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	en := en.New()
	uni = ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

}

func init() {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := field.Tag.Get("json")
		if name == "-" || name == "" {
			return field.Name // Fallback to struct field name if json tag is absent
		}
		return name
	})

	// custom validator
	validate.RegisterValidation("duration", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		_, err := time.ParseDuration(s)
		return err == nil
	})
}

func GetTranslator(locale string) ut.Translator {
	trans, _ := uni.GetTranslator("en")
	return trans
}

func ValidateStruct(obj any) error {
	return validate.Struct(obj)
}

func ValidateArray[T any](obj []T) error {
	var errs validator.ValidationErrors
	for _, v := range obj {
		err := ValidateStruct(v)
		if err != nil {
			var vErr validator.ValidationErrors
			if errors.As(err, &vErr) {
				errs = append(errs, vErr...)
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
