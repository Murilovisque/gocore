package gcvalidation

import (
	"cmp"
	"fmt"
	"reflect"
	"regexp"

	"github.com/Murilovisque/gocore/gctxt"
)

func EmailValidation(fieldName gctxt.Txt) Validation[string] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		switch lang {
		case gctxt.LangEnglish:
			return fmt.Sprintf("the field '%s' must be a valid e-mail address", txt)
		default:
			return fmt.Sprintf("campo '%s' deve ser um e-mail válido", txt)
		}
	})
	return func(vl string) (bool, gctxt.Txt) {
		if emailRegex.MatchString(vl) {
			return true, nil
		}
		return false, txt
	}
}

func GreaterThanValidation[T cmp.Ordered](fieldName gctxt.Txt, greaterThan T) Validation[T] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		switch lang {
		case gctxt.LangEnglish:
			return fmt.Sprintf("the field '%s' must be greater than '%v'", txt, greaterThan)
		default:
			return fmt.Sprintf("o campo '%s' deve ser maior que '%v'", txt, greaterThan)
		}
	})
	return func(vl T) (bool, gctxt.Txt) {
		if vl > greaterThan {
			return true, nil
		}
		return false, txt
	}
}

func GreaterOrEqualValidation[T cmp.Ordered](fieldName gctxt.Txt, greaterThan T) Validation[T] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		switch lang {
		case gctxt.LangEnglish:
			return fmt.Sprintf("the field '%s' must be greater than or equal to '%v'", txt, greaterThan)
		default:
			return fmt.Sprintf("o campo '%s' deve ser menor ou igual a '%v'", txt, greaterThan)
		}

	})
	return func(vl T) (bool, gctxt.Txt) {
		if vl >= greaterThan {
			return true, nil
		}
		return false, txt
	}
}

func LessThanValidation[T cmp.Ordered](fieldName gctxt.Txt, lessThan T) Validation[T] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		switch lang {
		case gctxt.LangEnglish:
			return fmt.Sprintf("the field '%s' must be less than '%v'", txt, lessThan)
		default:
			return fmt.Sprintf("o campo '%s' deve ser menor que '%v'", txt, lessThan)
		}
	})
	return func(vl T) (bool, gctxt.Txt) {
		if vl < lessThan {
			return true, nil
		}
		return false, txt
	}
}

func LessOrEqualValidation[T cmp.Ordered](fieldName gctxt.Txt, lessThan T) Validation[T] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		switch lang {
		case gctxt.LangEnglish:
			return fmt.Sprintf("the field '%s' must be less than or equal to '%v'", txt, lessThan)
		default:
			return fmt.Sprintf("o campo '%s' deve ser menor ou igual a '%v'", txt, lessThan)
		}

	})
	return func(vl T) (bool, gctxt.Txt) {
		if vl <= lessThan {
			return true, nil
		}
		return false, txt
	}
}

func BetweenInclusiveValidation[T cmp.Ordered](fieldName gctxt.Txt, greaterOrEqualThan T, lessOrEqualThan T) Validation[T] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		switch lang {
		case gctxt.LangEnglish:
			return fmt.Sprintf("the field '%s' must be greater than or equal to '%v' and less than or equal to '%v'", txt, greaterOrEqualThan, lessOrEqualThan)
		default:
			return fmt.Sprintf("o campo '%s' deve ser maior ou igual a '%v' e menor ou igual a '%v'", txt, greaterOrEqualThan, lessOrEqualThan)
		}
	})
	return func(vl T) (bool, gctxt.Txt) {
		if vl >= greaterOrEqualThan && vl <= lessOrEqualThan {
			return true, nil
		}
		return false, txt
	}
}

func BetweenValidation[T cmp.Ordered](fieldName gctxt.Txt, greaterThan T, lessThan T) Validation[T] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		switch lang {
		case gctxt.LangEnglish:
			return fmt.Sprintf("the field '%s' must be greater than '%v' and less than '%v'", txt, greaterThan, lessThan)
		default:
			return fmt.Sprintf("o campo '%s' deve ser maior que '%v' e menor que '%v'", txt, greaterThan, lessThan)
		}

	})
	return func(vl T) (bool, gctxt.Txt) {
		if vl > greaterThan && vl < lessThan {
			return true, nil
		}
		return false, txt
	}
}

func RequiredValidation[T comparable](fieldName gctxt.Txt) Validation[T] {
	var zeroValue T
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		switch lang {
		case gctxt.LangEnglish:
			return fmt.Sprintf("the field '%s' is required", txt)
		default:
			return fmt.Sprintf("o campo '%s' é obrigatório", txt)
		}
	})
	return func(vl T) (bool, gctxt.Txt) {
		if vl != zeroValue {
			return true, nil
		}
		return false, txt
	}
}

func NotNilValidation[T any](fieldName gctxt.Txt) Validation[T] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		switch lang {
		case gctxt.LangEnglish:
			return fmt.Sprintf("the field '%s' cannot be nil", txt)
		default:
			return fmt.Sprintf("o campo '%s' não pode ser nulo", txt)
		}
	})
	return func(vl T) (bool, gctxt.Txt) {
		v := reflect.ValueOf(vl)
		if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface || v.Kind() == reflect.Map || v.Kind() == reflect.Slice || v.Kind() == reflect.Chan || v.Kind() == reflect.Func {
			if v.IsNil() {
				return false, txt
			}
		}
		return true, nil
	}
}

func PatternValidation(fieldName gctxt.Txt, regex *regexp.Regexp, acceptPattern string) Validation[string] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		switch lang {
		case gctxt.LangEnglish:
			return fmt.Sprintf("the '%s' field accepts the text pattern '%s'.", txt, acceptPattern)
		default:
			return fmt.Sprintf("o campo '%s' aceita o padrão de texto '%s'", txt, acceptPattern)
		}
	})
	return func(vl string) (bool, gctxt.Txt) {
		return regex.MatchString(vl), txt
	}
}
