package framework

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func (b *Backend) Validate(in interface{}) error {
	err := validator.New().Struct(in)
	if err != nil {
		ferr := err.(validator.ValidationErrors)
		var fs []string
		for _, fieldError := range ferr {
			fs = append(fs, fieldError.Field())
		}
		return fmt.Errorf("缺少参数: %s", strings.Join(fs, ", "))
	}
	return nil
}
