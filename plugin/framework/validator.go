package framework

import "github.com/go-playground/validator/v10"

func (b *Backend) Validate(in interface{}) error {
	return validator.New().Struct(in)
}
