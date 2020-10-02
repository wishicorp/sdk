package logical

import (
	"errors"
	"github.com/wishicorp/sdk/helper/jsonutil"
)

const Authorization = "Authorization"
const AuthTokenName = "x-auth-token"

var ErrAuthMethodRequired = errors.New("Authorization method required")
var ErrAuthMethodNotFound = errors.New("Authorization method not found")
var ErrAuthorizationTokenRequired = errors.New("Authorization token required")
var ErrAuthorizationTokenInvalid = errors.New("Authorization token invalid")

//验证信息
type Authorized map[string]interface{}

func (a Authorized) Encode() ([]byte, error) {
	return jsonutil.EncodeJSON(a)
}

func (a Authorized) GetAuthorizer() interface{} {
	return a[Authorization]
}
func (a Authorized) SetAuthorizer(in interface{}) {
	a[Authorization] = in
}
