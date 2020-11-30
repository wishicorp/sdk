package framework

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func (b *Backend) Validate(in interface{}) error {
	err := validator.New().Struct(in)
	var fs []string
	if err != nil {
		switch ferr := err.(type) {
		case validator.ValidationErrors:
			for _, fieldError := range ferr {
				fs = append(fs, getTag(in, fieldError.StructField()))
			}
		default:
			fs = append(fs, err.Error())
		}
		return fmt.Errorf("缺少参数: %s", strings.Join(fs, ", "))
	}
	return nil
}
func getTag(in interface{}, name string)string  {
	obj := reflect.TypeOf(in)
	var field reflect.StructField
	switch obj.Kind() {
	case reflect.Ptr:
		f, ok := obj.Elem().FieldByName(name)
		if ok{field = f}
	case reflect.Struct:
		f, ok := obj.FieldByName(name)
		if ok{field = f}
	default:
		return ""
	}
	if field.Tag.Get("json")!=""{
		return field.Tag.Get("json")
	}else {
		return field.Name
	}
}