package entities

type baseErrors struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func (e *baseErrors) Error() string {
	return e.Message
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
