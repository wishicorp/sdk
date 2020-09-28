package framework

import (
	"fmt"
	"github.com/wishicorp/sdk/plugin/logical"
	"strings"
)

// Helper which returns a generic regex string for creating endpoint patterns
// that are identified by the given name in the backends
func GenericNameRegex(name string) string {
	return fmt.Sprintf("(?P<%s>\\w(([\\w-.]+)?\\w)?)", name)
}

// GenericNameWithAtRegex returns a generic regex that allows alphanumeric
// characters along with -, . and @.
func GenericNameWithAtRegex(name string) string {
	return fmt.Sprintf("(?P<%s>\\w(([\\w-.@]+)?\\w)?)", name)
}

// Helper which returns a regex string for optionally accepting the a field
// from the API URL
func OptionalParamRegex(name string) string {
	return fmt.Sprintf("(/(?P<%s>.+))?", name)
}

// Helper which returns a regex string for capturing an entire endpoint path
// as the given name.
func MatchAllRegex(name string) string {
	return fmt.Sprintf(`(?P<%s>.*)`, name)
}

// PathAppend is a helper for appending lists of paths into a single
// list.
func PathAppend(paths ...[]*Namespace) []*Namespace {
	result := make([]*Namespace, 0, 10)
	for _, ps := range paths {
		result = append(result, ps...)
	}

	return result
}

// Namespace is a single path that the grpc-backend responds to.
type Namespace struct {
	Pattern        string
	Description    string
	Operations     map[logical.Operation]OperationHandler
	ExistenceCheck ExistenceFunc
	Deprecated     bool
}

// OperationHandler defines and describes a specific operation handler.
type OperationHandler interface {
	Handler() OperationFunc
	Properties() OperationProperties
}

// OperationProperties describes an operation for documentation, help text,
// and other clients. A Summary should always be provided, whereas other
// fields can be populated as needed.
type OperationProperties struct {
	Description string
	Authorized  bool
	Deprecated  bool
	Input       *logical.SchemaType `json:"-"`
	Output      *logical.SchemaType `json:"-"`
}

type Response struct {
	Description string // summary of the the response and should always be provided
	MediaType   string // media type of the response, defaulting to "application/json" if empty
}

// PathOperation is a concrete implementation of OperationHandler.
type PathOperation struct {
	Callback    OperationFunc
	Description string
	Authorized  bool
	Deprecated  bool
	Input       *logical.SchemaType
	Output      *logical.SchemaType
}

func (p *PathOperation) Handler() OperationFunc {
	return p.Callback
}

func (p *PathOperation) Properties() OperationProperties {
	return OperationProperties{
		Description: strings.TrimSpace(p.Description),
		Deprecated:  p.Deprecated,
		Authorized:  p.Authorized,
		Input:       p.Input,
		Output:      p.Output,
	}
}
