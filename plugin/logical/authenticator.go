package logical

import (
	"errors"
	"github.com/wishicorp/sdk/helper/jsonutil"
)

const AuthTokenName = "x-auth-token"
const Authorization = "authorization"

var ErrAuthMethodRequired = errors.New("Authorization method required")
var ErrAuthMethodNotFound = errors.New("Authorization method not found")
var ErrAuthorizationTokenRequired = errors.New("Authorization token required")
var ErrAuthorizationTokenInvalid = errors.New("Authorization token invalid")

//验证信息
type Authorized struct {
	ID            interface{} `json:"id" name:"账户ID"`
	Authorization string      `json:"authorization" name:"authorization token"`
	Principal     Principal   `json:"principal" name:"账户凭证(用户信息)"`
}
type Principal map[string]interface{}

func NewAuthorized(id interface{}, authorization string, principal Principal) Authorized {
	return Authorized{ID: id, Authorization: authorization, Principal: principal}
}

func (a Authorized) Encode() ([]byte, error) {
	return jsonutil.EncodeJSON(a)
}

func (a Authorized) GetPrincipal() Principal {
	return a.Principal
}

func (a Authorized) SetPrincipal(in Principal) {
	a.Principal = in
}
