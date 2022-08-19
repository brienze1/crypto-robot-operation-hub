package exceptions

type HandlerError struct {
	Message         string `json:"message"`
	InternalMessage string `json:"exceptions"`
}

func (err HandlerError) Error() string {
	err.InternalMessage = "Error occurred while handling the event"
	return err.Message
}
