package httpapi

import (
	"atc/atc"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Decode[T any](r *http.Request) (T, error) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var v T
	err := dec.Decode(&v)
	if err != nil {
		var ute *json.UnmarshalTypeError
		if errors.As(err, &ute) {
			return v, atc.Error{
				Code:    atc.ErrCodeInvalid,
				Message: fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset),
			}
		}

		var se *json.SyntaxError
		if errors.As(err, &se) {
			return v, atc.Error{
				Code:    atc.ErrCodeInvalid,
				Message: fmt.Sprintf("Syntax error (offset=%v, error=%v)", se.Offset, se.Error()),
			}
		}

		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return v, atc.Error{
				Code:    atc.ErrCodeInvalid,
				Message: "Invalid input.",
			}
		}

		return v, atc.Error{
			Code:     atc.ErrCodeInternal,
			Message:  "Internal error.",
			Internal: err,
		}
	}
	return v, nil
}
