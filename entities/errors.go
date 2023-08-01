package entities

type baseErrors struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func (e *baseErrors) Error() string {
	return e.Message
}

func NewError(message string, details []string) error {
	return &baseErrors{
		Message: message,
		Details: details,
	}
}

type AccountAlreadyExistsError struct {
	*baseErrors
}

func NewAccountAlreadyExistsError(
	documentNumber string,
) *AccountAlreadyExistsError {
	return &AccountAlreadyExistsError{
		baseErrors: &baseErrors{
			Message: "Account already exists",
			Details: []string{documentNumber},
		},
	}
}
