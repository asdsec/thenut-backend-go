package api

import (
	"github.com/asdsec/thenut/utils"
	"github.com/go-playground/validator/v10"
)

var validGender validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if gender, ok := fieldLevel.Field().Interface().(string); ok {
		return utils.IsSupportedGender(gender)
	}
	return false
}
