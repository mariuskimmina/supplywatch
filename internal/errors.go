package internal
/*
import "fmt"

type Error struct {
	orig error
	msg  string
	code ErrorCode
}

type ErrorCode uint

const (
	ErrorCodeUnkown ErrorCode = iota
	ErrorCodeNotFund
	ErrorCodeInvalidArgument
)

func WrapErrorf(orig error, code ErrorCode, format string, a ...interface{}) error {
	return &Error{
		code: code,
		orig: orig,
		msg:  fmt.Sprint(format, a...),
	}
}

func NewErrorf(code ErrorCode, format string, a ...interface{}) error {
	return WrapErrorf(nil, code, format, a...)
}

func (e *Error) Error() string {
	if e.orig != nil {
		return fmt.Sprint("%s: %v", e.msg, e.orig)
	}
	return e.msg
}

func (e *Error) Unwrap() error {
	return e.orig
}

func (e *Error) Code() ErrorCode {
	return e.code
}
*/
