package logical

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/copystructure"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"strings"
)

var ErrInvalidData = errors.New("invalid request data")

type Request struct {
	ID         string              `json:"id" structs:"id" mapstructure:"id"`
	Operation  Operation           `json:"operation" structs:"operation" mapstructure:"operation"`
	Namespace  string              `json:"namespace" structs:"path" mapstructure:"namespace"`
	Data       map[string][]byte   `json:"map" structs:"data" mapstructure:"data"`
	Authorized *Authorized         `json:"authorized" structs:"authorized" mapstructure:"authorized"`
	Token      string              `json:"token" structs:"token" mapstructure:"token"`
	Headers    map[string][]string `json:"headers" structs:"headers" mapstructure:"headers"`
	Connection *Connection         `json:"connection" structs:"connection" mapstructure:"connection"`
}

func (r *Request) GetData() []byte {
	input, ok := r.Data["data"]
	if !ok {
		return nil
	}
	return input
}

func (r *Request) Decode(out interface{}) error {
	if r.Data == nil {
		return ErrInvalidData
	}
	input, ok := r.Data["data"]
	if !ok {
		return ErrInvalidData
	}
	err := jsonutil.DecodeJSON(input, out)
	if err != nil {
		return err
	}
	return r.Validate(out)
}

func (r *Request) XMLDecode(out interface{}) error {
	if r.Data == nil {
		return ErrInvalidData
	}
	input, ok := r.Data["data"]
	if !ok {
		return ErrInvalidData
	}
	err := xml.Unmarshal(input, out)
	if err != nil {
		return err
	}
	return r.Validate(out)
}
func (r *Request) Validate(in interface{}) error {
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
func (r *Request) GetAuthorized() *Authorized {
	return r.Authorized
}
func (r *Request) SetAuthorized(a *Authorized) {
	r.Authorized = a
}

// Clone returns a deep copy of the request by using copystructure
func (r *Request) Clone() (*Request, error) {
	cpy, err := copystructure.Copy(r)
	if err != nil {
		return nil, err
	}
	return cpy.(*Request), nil
}

// Backend returns a data field and guards for nil Data
func (r *Request) Get(key string) interface{} {
	if r.Data == nil {
		return nil
	}
	return r.Data[key]
}

// GetString returns a data field as a string
func (r *Request) GetString(key string) string {
	raw := r.Get(key)
	s, _ := raw.(string)
	return s
}

func (r *Request) GoString() string {
	return fmt.Sprintf("*%#v", *r)
}
func (r *Request) String() string {
	value, _ := jsonutil.EncodeJSON(r)
	return string(value)
}

// Operation is an enum that is used to specify the type
// of request being made
type Operation string

const (
	IndexOption Operation = "index"
)

// InitializationRequest stores the parameters and context of an Initialize()
// call being made to a logical.Backend.
type InitializationRequest struct {
	Params map[string]string
}
