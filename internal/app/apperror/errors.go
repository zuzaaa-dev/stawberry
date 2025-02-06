package apperror

import (
	"fmt"
)

const (
	NotFound       = "NOT_FOUND"
	DatabaseError  = "DATABASE_ERROR"
	InternalError  = "INTERNAL_ERROR"
	DuplicateError = "DUPLICATE_ERROR"
	BadRequest     = "BAD_REQUEST"
	Unauthorized   = "UNAUTHORIZED"
	InvalidToken   = "INVALID_TOKEN"
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

type OfferError struct {
	Code    string
	Message string
	Err     error
}

func (e *OfferError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

var ErrOfferNotFound = &OfferError{
	Code:    NotFound,
	Message: "product not found",
}

type UserError struct {
	Code    string
	Message string
	Err     error
}

func (e *UserError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

var (
	ErrUserNotFound = &ProductError{
		Code:    NotFound,
		Message: "product not found",
	}
	ErrIncorrectPassword = &ProductError{
		Code:    Unauthorized,
		Message: "incorrect password",
	}
	ErrFailedToGeneratePassword = &ProductError{
		Code:    InternalError,
		Message: "failed to generate password",
	}
)

type TokenError struct {
	Code    string
	Message string
	Err     error
}

func (e *TokenError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

var ErrInvalidToken = &TokenError{
	Code:    InvalidToken,
	Message: "invalid token",
}
