package framework

import (
	"fmt"
	"github.com/wishicorp/sdk/plugin/logical"
)

func (b *Backend) initSchemaOnce() {
	_ = b.initSchema()
}

func (b *Backend) initSchema() error {
	schemas := logical.NamespaceSchemas{}
	for _, ns := range b.Namespaces {
		if ns.Description == "" {
			return fmt.Errorf("namespace[%s] description required", ns.Pattern)
		}
		namespace := logical.NamespaceSchema{
			Namespace:   ns.Pattern,
			Description: ns.Description,
			Operations:  make(map[logical.Operation]*logical.Schema),
		}
		for opt, handler := range ns.Operations {
			properties := handler.Properties()
			if properties.Description == "" {
				return descriptionError(ns.Pattern, opt)
			}
			input, err := properties.Input.Fields()
			if err != nil {
				return schemaError(ns.Pattern, opt, err)
			}
			output, err := properties.Output.Fields()
			if err != nil {
				return schemaError(ns.Pattern, opt, err)
			}
			schema := &logical.Schema{
				Description: properties.Description,
				Authorized:  properties.Authorized,
				Deprecated:  properties.Deprecated,
				Input:       input,
				Output:      output,
			}
			namespace.Operations[opt] = schema
		}
		schemas = append(schemas, &namespace)
	}
	b.schemas = schemas
	return nil
}

func schemaError(pattern string, operation logical.Operation, err error) error {
	return fmt.Errorf("namespace[%s] operation[%s] %s", pattern, operation, err)
}
func descriptionError(pattern string, operation logical.Operation) error {
	return fmt.Errorf("namespace[%s] operation[%s] Description required", pattern, operation)
}
