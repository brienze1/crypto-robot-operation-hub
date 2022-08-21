package exceptions

type handlerError struct {
	message       string
	internalError string
	description   string
}

func HandlerError(err error, internalError string) *handlerError {
	return &handlerError{
		message:       err.Error(),
		internalError: internalError,
		description:   "Error occurred while handling the event",
	}
}

func (h *handlerError) Error() string {
	return h.message
}

func (h *handlerError) InternalError() string {
	return h.internalError
}

func (h *handlerError) Description() string {
	return h.description
}
