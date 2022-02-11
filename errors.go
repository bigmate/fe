package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FieldErrorf returns a field specific error
// that is useful to return structured error
func FieldErrorf(field, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	s := status.New(codes.InvalidArgument, msg)
	st, err := s.WithDetails(&FieldError{
		Field: field,
		Msg:   msg,
	})

	if err != nil {
		return err
	}

	return st.Err()
}

// NonFieldErrorf returns an error that is not related to any field
func NonFieldErrorf(format string, args ...interface{}) error {
	return FieldErrorf("nonField", format, args...)
}
