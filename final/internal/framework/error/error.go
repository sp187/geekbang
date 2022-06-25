package errs

import "net/http"

type Error struct {
	HttpCode int    `json:"-"`    // standard http code
	Code     int    `json:"code"` // app error code
	Msg      string `json:"msg"`  // error message
}

func (e *Error) Error() string {
	return e.Msg
}

var (
	Server     = &Error{Code: 601, Msg: "服务器异常,请稍后重试", HttpCode: http.StatusOK}
	Timeout    = &Error{Code: 602, Msg: "服务响应超时", HttpCode: http.StatusOK}
	BadRequest = &Error{Code: 603, Msg: "参数异常", HttpCode: http.StatusBadRequest}
)
