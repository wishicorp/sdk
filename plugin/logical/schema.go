package logical

import (
	"fmt"
	"reflect"
	"strings"
)

const SchemaRequestTag = "__SCHEMA__"

//接口属性列
type Field struct {
	Field      string   `json:"field"`
	Name       string   `json:"name"`
	Kind       string   `json:"kind"`
	Required   bool     `json:"required"`
	Deprecated bool     `json:"deprecated"`
	Reference  []*Field `json:"reference"`
	Example    string   `json:"example"`
	IsList     bool     `json:"is_list"`
}

type EmptySchema struct{}

type NamespaceSchemas []*NamespaceSchema
type NamespaceSchema struct {
	Namespace   string `json:"namespace"`
	Description string `json:"description"`
	Operations  map[Operation]*Schema
}

type SchemaResponse struct {
	NamespaceSchemas []*NamespaceSchema `json:"namespace_schemas"`
}

//接口属性
type Schema struct {
	Description string   `json:"description"`
	Authorized  bool     `json:"authorized"`
	Deprecated  bool     `json:"deprecated"`
	Input       []*Field `json:"input,omitempty"`
	Output      []*Field `json:"output,omitempty"`
}

func fieldError(sName, fName, tName string) error {
	return fmt.Errorf("struct[%s] filed[%s] tag[%s] required", sName, fName, tName)
}

//列属性猜解
type SchemaType struct {
	Type reflect.Type
}

func getType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Ptr:
		return t.Elem()
	case reflect.Struct:
		return t
	case reflect.Slice:
		fallthrough
	case reflect.Map:
		return getType(t.Elem())
	default:
		return t
	}
}

func getFields(Type reflect.Type) []*Field {
	defer func() {
		recover()
	}()
	var fields []*Field
	for i := 0; i < Type.NumField(); i++ {
		field := new(Field)
		f := Type.Field(i)
		isList := f.Type.Kind() == reflect.Slice
		realType := getType(f.Type)
		kindString := realType.Kind().String()
		//fmt.Println(f.Name, realType.Kind(), realType.Kind().String(), realType.Name())
		if realType.Kind() == reflect.Struct && realType.Name() != "Time" {
			field.Reference = getFields(realType)
		}
		if realType.Name() == "Time" {
			kindString = "datetime"
		}

		fValue := f.Tag.Get("json")
		if fValue == "" {
			fValue = f.Name
		}
		fName := f.Tag.Get("name")
		if fName == "" {
			fName = f.Name
		}

		example := f.Tag.Get("example")
		validate := f.Tag.Get("validate")
		required := validate != "" && strings.Contains(strings.ToLower(validate), "required")

		if strings.Contains(strings.ToUpper(f.Tag.Get("xorm")), "NOT NULL") {
			required = true
		}

		field.Name = fName
		field.Field = fValue
		field.Required = required
		field.Kind = kindString
		field.Deprecated = f.Tag.Get("deprecated") != ""
		field.Example = example
		field.IsList = isList
		fields = append(fields, field)
	}
	return fields
}
func (s *SchemaType) Fields() ([]*Field, error) {
	return getFields(s.Type), nil
}
