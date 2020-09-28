package consts

type ReplyCode int

const (
	ReplyCodeSuccess            ReplyCode = 0
	ReplyCodeFailure            ReplyCode = 1
	ReplyCodeMethodInvalid      ReplyCode = 2
	ReplyCodeSignInvalid        ReplyCode = 3
	ReplyCodeReqBlocked         ReplyCode = 4
	ReplyCodeRateLimited        ReplyCode = 5
	ReplyCodeBackendNotExists   ReplyCode = 6
	ReplyCodeBackendShutdown    ReplyCode = 7
	ReplyCodeNamespaceNotExists ReplyCode = 8
	ReplyCodeOperationNotExists ReplyCode = 9
	ReplyCodeAuthorizedRequired ReplyCode = 10
)

func (r ReplyCode) Code() int {
	return int(r)
}
