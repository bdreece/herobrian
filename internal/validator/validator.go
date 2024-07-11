package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
    inner *validator.Validate
}

var Default = &Validator{validator.New()}

func (v *Validator) Validate(i any) error {
    return v.inner.Struct(i)
}
