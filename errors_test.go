package fe

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fieldError interface {
	GetField() string
	GetError() string
}

func TestFieldErrorf(t *testing.T) {
	err := Field("name", "invalid name").
		Field("password", "invalid password").
		Field("avatar", "pass base64 encoded string").
		NonField("unexpected error").
		Field("name", "too long").
		Error()

	st, ok := status.FromError(err)

	if !ok {
		t.Fatal("expected status")
	}

	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected %v, got %v", codes.InvalidArgument, st.Code())
	}

	for _, detail := range st.Details() {
		fe, ok := detail.(fieldError)

		if !ok {
			t.Fatal("expected field error")
		}

		t.Log(fe.GetField(), fe.GetError())
	}
}
