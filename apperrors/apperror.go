package apperrors

import "fmt"

const (
	_         = iota
	E_FATAL   // panic error
	E_READ    // error reading a file
	E_CREATE  // error crating a file
	E_API     // error api
	E_PARSE   // parse error or convert
	E_REQUIRE // error required data
	E_CONF    // error configuration
)

type AppError struct {
	Msg      string // message about error
	NameFunc string // name of function where was a error
	Status   int    // status of error
	Err      error
}

func (a *AppError) Error() string {
	return fmt.Sprintf("Msg: %s, Function: %s, Status: %d, err: %s", a.Msg, a.NameFunc, a.Status, a.Err.Error())
}
