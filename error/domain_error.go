package error

import "fmt"

type DomainError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error { return e.Err }

func NewError1(code int, message string) *DomainError{
	return &DomainError{Code: code, Message: message}
}

func NewError2(code int, message string, err error) *DomainError {
	return &DomainError{Code: code, Message: message, Err: err}
}
