package validationerrors

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	name     string
	messages []string
}

func (err *ValidationError) AddMessage(message string) {
	err.messages = append(err.messages, message)
}

func (err ValidationError) Error() string {
	return fmt.Sprintf(
		"%s invalid: %s",
		err.name,
		strings.Join(err.messages, ", "),
	)
}

func (err *ValidationError) Present() bool {
	return len(err.messages) > 0
}

func NewValidationError(name string) *ValidationError {
	return &ValidationError{
		name: name,
	}
}
