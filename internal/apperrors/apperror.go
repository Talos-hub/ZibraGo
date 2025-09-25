package apperrors

import "fmt"

const (
	_         = iota
	E_FATAL   // panic error
	E_READ    // error reading a file
	E_CREATE  // error create
	E_API     // error api
	E_PARSE   // parse error or convert
	E_REQUIRE // error required data
	E_CONF    // error configuration
	E_OPEN
	E_MULTIPLE // when got a lot of errors
	E_ADD
	E_EMPTY_DIR
)

type AppError struct {
	Msg      string // message about error
	NameFunc string // name of function where was a error
	Status   int    // status of error
	Err      error
}

func NewAppError(msg, nameFunc string, status int, err error) *AppError {
	return &AppError{
		Msg:      msg,
		NameFunc: nameFunc,
		Status:   status,
		Err:      err,
	}
}
func (a *AppError) Error() string {
	return fmt.Sprintf("Msg: %s, Function: %s, Status: %d, err: %s", a.Msg, a.NameFunc, a.Status, a.Err.Error())
}
