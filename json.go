package fe

import (
	"bytes"
	"encoding/json"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func JSON(err error) (io.Reader, bool) {
	s, ok := status.FromError(err)
	if !ok || s.Code() != codes.InvalidArgument {
		return nil, false
	}

	bucket := make(map[string][]string)

	for _, detail := range s.Details() {
		if det, fok := detail.(*FieldError); fok {
			bucket[det.GetField()] = append(bucket[det.GetField()], det.GetError())
		}
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(bucket)

	return buf, true
}
