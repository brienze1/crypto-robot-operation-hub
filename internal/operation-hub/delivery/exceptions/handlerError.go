package exceptions

type HandlerError struct {
	Message         string `json:"message"`
	InternalMessage string `json:"exceptions"`
	Description     string `json:"description"`
}

func (err *HandlerError) Error() string {
	err.Description = "Error occurred while handling the event"
	return err.Message
}
