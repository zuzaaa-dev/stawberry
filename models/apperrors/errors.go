package apperrors

import (
	"fmt"
)

const (
	NotFound       = "NOT_FOUND"
	DatabaseError  = "DATABASE_ERROR"
	InternalError  = "INTERNAL_ERROR"
	DuplicateError = "DUPLICATE_ERROR"
)

type ProductError struct {
	Code    string
	Message string
	Err     error
}

func (e *ProductError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

var (
	ErrProductNotFound = &ProductError{
		Code:    NotFound,
		Message: "product not found",
	}
	ErrStoreNotFound = &ProductError{
		Code:    NotFound,
		Message: "store not found",
	}
)
