package gateway

type Security interface {
	SignVerify(args *RequestArgs) bool
	RateLimiter(method *Method, client *Client) error
	Blocker(method *Method, client *Client) error
}
