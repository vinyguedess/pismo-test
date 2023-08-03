package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		description          string
		expectedError        error
		expectedErrorMessage string
		expectedErrorType    error
	}{
		{
			description:          "Base Error",
			expectedErrorMessage: "Base error",
			expectedError:        NewError("Base error", []string{}),
			expectedErrorType:    &baseErrors{},
		},
		{
			description:          "AccountAlreadyExistsError",
			expectedErrorMessage: "Account already exists",
			expectedError:        NewAccountAlreadyExistsError("123456789"),
			expectedErrorType:    &AccountAlreadyExistsError{},
		},
		{
			description:          "ItemNotFoundError",
			expectedErrorMessage: "Item not found",
			expectedError:        NewItemNotFoundError("Item", 1),
			expectedErrorType:    &ItemNotFoundError{},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			assert.Equal(t, test.expectedErrorMessage, test.expectedError.Error())
			assert.IsType(t, test.expectedErrorType, test.expectedError)
		})
	}
}
