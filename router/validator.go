package router

import (
	"gopkg.in/go-playground/validator.v9"
)

type MyValidator struct {
	validator *validator.Validate
}

func (cv *MyValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// func validateMdn(fl validator.FieldLevel) bool {
// 	mdn := fl.Field().String()
// 	if b, err := regexp.MatchString(`(088|6288|88)[0-9]*$`, mdn); err != nil {
// 		return false
// 	} else {
// 		return b
// 	}
// 	//return (len(mdn) >= 8 && len(mdn) <= 15)
// }
