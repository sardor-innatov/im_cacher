package protocol

import "errors"

var (
	ErrEmptyRequest   = errors.New("error: empty request")
	ErrInvalidFormat  = errors.New("error: invalid request format, unknown command")
	ErrUnknownCommand = errors.New("error: unknown command")
	ErrInvalidSet     = errors.New("error: set command payload must be key:value")
)
