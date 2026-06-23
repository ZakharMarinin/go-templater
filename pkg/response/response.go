package response

import "errors"

var (
	ErrCanceled = errors.New("operation canceled")
	ErrNotExist = errors.New("go.mod not exist")
)
