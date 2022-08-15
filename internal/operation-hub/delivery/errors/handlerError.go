package errors

type HandlerError struct {
	Message         string `json:"message"`
	InternalMessage string `json:"errors"`
}

func (err HandlerError) Error() string {
	err.InternalMessage = "Error occurred while handling the event"
	return err.Message
}
