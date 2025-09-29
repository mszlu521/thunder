package errs

var ErrParam = NewError(400, "param error")
var ErrUnauthorized = NewError(401, "unauthorized")
var NoEventHandler = NewError(500, "no handler")
var DBError = NewError(999, "db error")

type Errors struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *Errors) Error() string {
	return e.Msg
}

func NewError(code int, msg string) *Errors {
	return &Errors{
		Code: code,
		Msg:  msg,
	}
}
