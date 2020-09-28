package httputil

// http返回分页数据结构
type pagination struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// http返回数据结构
type Response struct {
	Code       int         `json:"code"`
	Result     interface{} `json:"result,omitempty"`
	Message    string      `json:"message,omitempty"`
	Pagination *pagination `json:"pagination,omitempty"`
}

func Pagination(total int, page *Pageable) *pagination {
	mod := total % page.Size()
	totalPages := total / page.Size()
	if 0 == totalPages {
		totalPages = 1
	}
	if mod > 0 && total > page.Size() {
		totalPages += 1
	}
	return &pagination{Total: total, TotalPages: totalPages, Page: page.Page(), Size: page.Size()}
}
func PageResponse(pagination *pagination, data interface{}) *Response {
	return &Response{Code: 0, Pagination: pagination, Result: data}
}

func ErrorResponse(code int, message string) *Response {
	return &Response{Code: code, Message: message}
}

func SuccessResponse(data interface{}) *Response {
	return &Response{Code: 0, Result: data}
}
