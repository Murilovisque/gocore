package gcvalidation

import (
	"cmp"
	"fmt"
	"regexp"

	"github.com/Murilovisque/gocore/gctxt"
)

func EmailValidation(fieldName gctxt.Txt) Validation[string] {
	emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		return fmt.Sprintf("campo '%s' deve ser um e-mail válido", txt)
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
		return fmt.Sprintf("campo '%s' deve ser menor que '%v'", txt, greaterThan)
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
		return fmt.Sprintf("campo '%s' deve ser menor ou igual a '%v'", txt, greaterThan)
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
		return fmt.Sprintf("campo '%s' deve ser menor que '%v'", txt, lessThan)
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
		return fmt.Sprintf("campo '%s' deve ser menor ou igual a '%v'", txt, lessThan)
	})
	return func(vl T) (bool, gctxt.Txt) {
		if vl <= lessThan {
			return true, nil
		}
		return false, txt
	}
}

func BetweenInclusiveValidation[T cmp.Ordered](fieldName gctxt.Txt, greaterOrEqualThan T, lessOrEqualThan T) ValidationTwoFields[T, T] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		return fmt.Sprintf("campo '%s' deve ser maior ou igual a '%v' e menor ou igual a '%v'", txt, greaterOrEqualThan, lessOrEqualThan)
	})
	return func(vl1, vl2 T) (bool, gctxt.Txt) {
		if vl1 >= greaterOrEqualThan && vl2 <= lessOrEqualThan {
			return true, nil
		}
		return false, txt
	}
}

func BetweenValidation[T cmp.Ordered](fieldName gctxt.Txt, greaterThan T, lessThan T) ValidationTwoFields[T, T] {
	txt := fieldName.MapContent(func(lang gctxt.Language, txt string) string {
		return fmt.Sprintf("campo '%s' deve ser maior que '%v' e menor que '%v'", txt, greaterThan, lessThan)
	})
	return func(vl1, vl2 T) (bool, gctxt.Txt) {
		if vl1 > greaterThan && vl2 < lessThan {
			return true, nil
		}
		return false, txt
	}
}
