package fe

import (
	"encoding/json"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test(t *testing.T) {
	messages := map[string][]string{
		"name":     {"invalid name", "too long"},
		"password": {"invalid password"},
		"avatar":   {"pass base64 encoded string"},
		"nonField": {"unexpected error"},
	}
	fe := &chain{status.New(codes.InvalidArgument, "field error")}

	for field, msg := range messages {
		if field == "nonField" {
			for _, s := range msg {
				fe.NonField(s)
			}
		} else {
			for _, s := range msg {
				fe.Field(field, s)
			}
		}
	}

	err := fe.Error()

	t.Run("Chainer", func(t *testing.T) {
		st, ok := status.FromError(err)

		if !ok {
			t.Fatal("expected status")
		}

		if st.Code() != codes.InvalidArgument {
			t.Fatalf("expected %v, got %v", codes.InvalidArgument, st.Code())
		}

		for _, detail := range st.Details() {
			fe, ok := detail.(*FieldError)

			if !ok {
				t.Fatal("expected field error")
			}

			t.Log(fe.GetField(), fe.GetError())
		}
	})

	t.Run("JSON", func(t *testing.T) {
		r, ok := JSON(err)
		if !ok {
			t.Fatalf("expected FieldError")
		}

		dest := make(map[string][]string)

		if decErr := json.NewDecoder(r).Decode(&dest); decErr != nil {
			t.Fatalf("unexpected error: %v", decErr)
		}

		if !reflect.DeepEqual(dest, messages) {
			t.Errorf("messages are not deeply equal: %v\n %v\n", dest, messages)
		}
	})
}
