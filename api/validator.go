package api

import (
	"backend-in-go/util"
	"github.com/go-playground/validator/v10"
)

var validatorCurrency validator.Func = func(level validator.FieldLevel) bool {
	if currency, ok := level.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
