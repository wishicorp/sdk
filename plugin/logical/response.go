package logical

// Response is a struct that stores the response of a request.
// It is used to abstract the details of the higher level request protocol.
type Response struct {
	ResultCode int64               `json:"result_code" structs:"result_code" mapstructure:"result_code"`
	ResultMsg  string              `json:"result_msg,omitempty" structs:"result_msg" mapstructure:"result_msg"`
	Data       interface{}         `json:"data,omitempty" structs:"data" mapstructure:"data"`
	Headers    map[string][]string `json:"headers,omitempty" structs:"headers" mapstructure:"headers"`
}
