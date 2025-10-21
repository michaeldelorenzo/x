package validationerrors_test

import (
	"testing"

	"github.com/michaeldelorenzo/x/v2/pkg/errors/validationerrors"
)

func TestValidationError(t *testing.T) {
	tests := []struct {
		name             string
		validationError  *validationerrors.ValidationError
		messagesToAdd    []string
		expectedErrorMsg string
		expectedPresent  bool
	}{
		{
			name:             "No messages",
			validationError:  validationerrors.NewValidationError("test"),
			messagesToAdd:    []string{},
			expectedErrorMsg: "test invalid: ",
			expectedPresent:  false,
		},
		{
			name:             "One message",
			validationError:  validationerrors.NewValidationError("test"),
			messagesToAdd:    []string{"message 1"},
			expectedErrorMsg: "test invalid: message 1",
			expectedPresent:  true,
		},
		{
			name:             "Multiple messages",
			validationError:  validationerrors.NewValidationError("test"),
			messagesToAdd:    []string{"message 1", "message 2"},
			expectedErrorMsg: "test invalid: message 1, message 2",
			expectedPresent:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.messagesToAdd {
				tt.validationError.AddMessage(msg)
			}

			if tt.validationError.Error() != tt.expectedErrorMsg {
				t.Errorf("Expected error message '%s', but got '%s'", tt.expectedErrorMsg, tt.validationError.Error())
			}

			if tt.validationError.Present() != tt.expectedPresent {
				t.Errorf("Expected present to be %t, but got %t", tt.expectedPresent, tt.validationError.Present())
			}
		})
	}
}
