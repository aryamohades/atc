package atc

import (
	"fmt"
)

const (
	ErrCodeInternal = "internal"
	ErrCodeInvalid  = "invalid"
)

type Error struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Details  any    `json:"details,omitempty"`
	Internal error  `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code=%s message=%s details=%v internal=%v", e.Code, e.Message, e.Details, e.Internal)
}
