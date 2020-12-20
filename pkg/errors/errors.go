package errors

import (
	"fmt"
)

type Code int

const (
	OtherClause    Code = -1
	InternalServer      = iota + 1
	DataNotFound
	DuplicateData
	InvalidQuery
	InvalidClassName
	InvalidFieldName
	ChangedImmutableField
	MissingRequiredField
	IncorrectFieldType
	InvalidJSON
	IncorrectOperation
)

type CustomError interface {
	Error() string
	Code() Code
}

type customError struct {
	code    Code
	message string
}

func New(code Code, message string) CustomError {
	return &customError{code, message}
}

func (e customError) Error() string {
	return fmt.Sprintf("error %d: %s", e.code, e.message)
}

func (e customError) Code() Code {
	return e.code
}
