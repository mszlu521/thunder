package errs

var ErrParam = NewError(400, "param error")
var NoEventHandler = NewError(500, "no handler")

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
