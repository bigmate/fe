package fe

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const nonField = "nonField"

//Chainer collects error for different fields and when Error() is called returns gRPC error
type Chainer interface {
	Field(field, format string, args ...interface{}) Chainer
	NonField(format string, args ...interface{}) Chainer
	Error() error
}

// Field returns a field specific error
// that is useful to return structured error
func Field(field, format string, args ...interface{}) Chainer {
	return (chain{status.New(codes.InvalidArgument, "field error")}).Field(field, format, args...)
}

// NonField returns an error that is not related to any field
func NonField(format string, args ...interface{}) Chainer {
	return Field(nonField, format, args...)
}

type chain struct {
	*status.Status
}

func (c chain) Field(field, format string, args ...interface{}) Chainer {
	st, _ := c.WithDetails(&FieldError{
		Field: field,
		Error: fmt.Sprintf(format, args...),
	})

	c.Status = st

	return c
}

func (c chain) NonField(format string, args ...interface{}) Chainer {
	return c.Field(nonField, format, args...)
}

func (c chain) Error() error {
	return c.Err()
}
