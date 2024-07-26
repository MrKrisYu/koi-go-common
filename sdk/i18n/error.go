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

func NewMyErrorWithMessageID(mi MessageID, arg ...any) MyError {
	myError := MyError{
		Err: nil,
		Message: Message{
			ID:             mi.ID,
			DefaultMessage: mi.DefaultMessage,
		},
	}
	if len(arg) > 0 {
		myError.Message.Args = arg[0]
	}
	return myError
}
