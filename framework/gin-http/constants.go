package gin_http

type ResponseCode int

const (
	RequestSuccess ResponseCode = 0 //操作成功
	SystemFail     ResponseCode = 1 //系统失败
	SystemFatal    ResponseCode = 2 //系统致命错误

	RequestFailed     ResponseCode = 1001 //请求失败
	RequestInvalid    ResponseCode = 1002 //请求无效
	RequestNotAllowed ResponseCode = 1003 //请求不允许
	RequestTimeout    ResponseCode = 1004 //请求超时

	DataBindFailed   ResponseCode = 2001 //数据绑定失败
	ParamsInvalid    ResponseCode = 2002 //参数无效
	ParamsBindFailed ResponseCode = 2003 //参数绑定失败

	DataNotFound     ResponseCode = 2011 //数据未找到
	DataUpdateFailed ResponseCode = 2012 //数据更新失败
	DataInsertFailed ResponseCode = 2013 //数据插入失败
	DataDeleteFailed ResponseCode = 2014 //数据删除失败

	InternalSvcFailed  ResponseCode = 3001 //内部服务失败
	InternalSvcTimeout ResponseCode = 3002 //内部服务超时

	ExternalSvcFailed      ResponseCode = 3011 //外部服务失败
	ExternalServiceTimeout ResponseCode = 3012 //外部服务超时

)

type CCodeMsg struct {
	CodeMsg map[ResponseCode]string
}

var CodeMsg = map[ResponseCode]string{
	RequestSuccess:         "操作成功",
	SystemFail:             "系统失败",
	SystemFatal:            "系统致命错误",
	RequestFailed:          "请求失败",
	RequestInvalid:         "请求无效",
	RequestNotAllowed:      "请求不允许",
	RequestTimeout:         "请求超时",
	DataBindFailed:         "数据绑定失败",
	ParamsInvalid:          "参数无效",
	ParamsBindFailed:       "参数绑定失败",
	DataNotFound:           "数据未找到",
	DataUpdateFailed:       "数据更新失败",
	DataInsertFailed:       "数据插入失败",
	DataDeleteFailed:       "数据删除失败",
	InternalSvcFailed:      "内部服务失败",
	InternalSvcTimeout:     "内部服务超时",
	ExternalSvcFailed:      "外部服务失败",
	ExternalServiceTimeout: "外部服务超时",
}

const (
	DefaultPageSize   int = 50
	DefaultPageNumber int = 1
)

const (
	Asc  string = "ASC"
	Desc string = "DESC"
)

func (code ResponseCode) CodeMsg() (ResponseCode, string) {
	return code, CodeMsg[code]
}
