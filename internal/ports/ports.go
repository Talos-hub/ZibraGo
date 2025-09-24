package ports

import "os"

// Walker walk to dir and returns pathes to files
type Walker interface {
	Walk(path string) ([]string, error)
}

// Archiver create new archiv and retrun it
type Archiver interface {
	Start(paths ...string) (string, error)
}

type Logger interface {
	Info(msg string, arg ...any)
	Warn(msg string, arg ...any)
	Debug(msg string, arg ...any)
	Error(msg string, arg ...any)
}

type ApiCloud interface {
	Check() error
	Auth() error
	SendFile(os *os.File) error
}
