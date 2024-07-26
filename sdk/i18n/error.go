package i18n

type MyError struct {
	Err     error
	Message Message `json:"message"`
}

func (e MyError) Error() string {
	return e.Err.Error()
}

func NewMyError(err error) MyError {
	return MyError{
		Err:     err,
		Message: Message{DefaultMessage: err.Error()},
	}
}
