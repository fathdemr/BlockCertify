package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
	Err     error
}

var (
	ErrInvalidFile         = "INVALID_FILE"
	ErrHashingFailed       = "HASHING_FAILED"
	ErrArweaveUploadFailed = "ARWEAVE_UPLOAD_FAILED"
	ErrBlockchainFailed    = "BLOCKCHAIN_FAILED"
	ErrDiplomaExists       = "DIPLOMA_EXISTS"
	ErrInsufficientBalance = "INSUFFICIENT_BALANCE"
	ErrVerificationFailed  = "VERIFICATION_FAILED"
	ErrTokenExpired        = "TOKEN_EXPIRED"
	ErrUserExists          = "USER_EXISTS"
	ErrInvalidCredentials  = "INVALID_CREDENTIALS"
	ErrTokenCreateFailed   = "TOKEN_CREATE_FAILED"
	ErrInvalidToken        = "INVALID_TOKEN"
	ErrInvalidRequest      = "INVALID_REQUEST"
)

func New(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}
