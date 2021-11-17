package internal

import (
	"fmt"
	"strings"
)

type Error struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Message string   `json:"message"`
	Trace   []string `json:"trace,omitempty"`
}

func NewError(value error) *Error {
	var errTrace string
	errTrace = fmt.Sprintf("%+v", value)
	errTrace = strings.ReplaceAll(errTrace, "\t", "    ")

	return &Error{
		Error: ErrorBody{
			Message: value.Error(),
			Trace:   strings.Split(errTrace, "\n"),
		},
	}
}
