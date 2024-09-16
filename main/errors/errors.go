package errors

import (
	"XrayHelper/main/serial"
	"reflect"
	"strings"
)

// Error is an error object with underlying error.
type Error struct {
	prefix  []any
	pathObj any
	message []any
}

// WithPrefix set err prefix in method Error()
func (err *Error) WithPrefix(prefix ...any) *Error {
	err.prefix = prefix
	return err
}

// WithPathObj set Obj path, should not be predeclared type like pointer, bool ...
func (err *Error) WithPathObj(obj any) *Error {
	err.pathObj = obj
	return err
}

func (err *Error) pkgPath() string {
	if err.pathObj == nil {
		return ""
	}
	path := reflect.TypeOf(err.pathObj).PkgPath()
	return path
}

// Error implements error.Error().
func (err *Error) Error() string {
	builder := strings.Builder{}
	for _, prefix := range err.prefix {
		builder.WriteByte('[')
		builder.WriteString(serial.ToString(prefix))
		builder.WriteString("] ")
	}

	path := err.pkgPath()
	if len(path) > 0 {
		builder.WriteString(path)
		builder.WriteString(": ")
	}

	msg := serial.Concat(err.message...)
	builder.WriteString(msg)

	return builder.String()
}

// New returns a new error object with message formed from given arguments
func New(msg ...any) *Error {
	return &Error{message: msg}
}
