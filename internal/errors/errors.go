package errors

import "errors"

var (
	ErrNoDBResult             = errors.New("no result from select in DB")
	ErrDuplicateValue         = errors.New("duplicate value while insert")
	ErrNoInfo                 = errors.New("no info")
	ErrIncorrectLoginPassword = errors.New("incorrect username or password")
)
