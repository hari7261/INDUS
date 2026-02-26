package cli

import (
	"context"
	"errors"
)

var (
	ErrMissingCommand = errors.New("missing command")
	ErrUnknownCommand = errors.New("unknown command")
	ErrInvalidInput   = errors.New("invalid input")
	ErrInternal       = errors.New("internal error")
)

type UserError struct {
	Msg string
}

func (e *UserError) Error() string {
	return e.Msg
}

type InternalError struct {
	Msg string
	Err error
}

func (e *InternalError) Error() string {
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	if errors.Is(err, context.Canceled) {
		return 130 // Standard exit code for SIGINT
	}

	var userErr *UserError
	if errors.As(err, &userErr) {
		return 2
	}

	if errors.Is(err, ErrMissingCommand) || errors.Is(err, ErrUnknownCommand) {
		return 2
	}

	return 1
}
