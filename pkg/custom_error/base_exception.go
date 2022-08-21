package custom_error

type BaseErrorAdapter interface {
	Error() string
	Description() string
	InternalError() string
}

type baseError struct {
	Message            string `json:"error"`
	InternalMessage    string `json:"internal_error"`
	DescriptionMessage string `json:"description"`
}

func BaseError(err error) *baseError {
	switch e := err.(type) {
	case BaseErrorAdapter:
		return &baseError{
			Message:            e.Error(),
			InternalMessage:    e.InternalError(),
			DescriptionMessage: e.Description(),
		}
	default:
		return &baseError{
			Message:            err.Error(),
			InternalMessage:    "",
			DescriptionMessage: "",
		}
	}
}
