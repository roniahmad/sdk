package sdk

import "errors"

const (
	Bearer                = "bearer"
	CreateDirectoryFailed = "Could not create user directory"
	ExpiredToken          = "Token is expired"
	InvalidAuthHeader     = "Invalid authorization header format"
	InvalidToken          = "Token is invalid"
	UnsupportedDocument   = "Unsupported document type"
)

var (
	ErrCreateDirectory     = newError(CreateDirectoryFailed)
	ErrExpiredToken        = newError(ExpiredToken)
	ErrInvalidAuthHeader   = newError(InvalidAuthHeader)
	ErrInvalidToken        = newError(InvalidToken)
	ErrUnsupportedDocument = newError(UnsupportedDocument)
)

func newError(message string) error {
	return errors.New(message)
}
