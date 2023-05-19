package api

import (
	"github.com/gen4ralz/simplebank/utils"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is support
		return utils.IsSupportCurrency(currency)
	}
	return false
}