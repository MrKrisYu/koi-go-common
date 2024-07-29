package i18n

type MyError struct {
	Err     error
	Message Message `json:"message"`
}

func (e MyError) Error() string {
	return e.Err.Error()
}

func NewMyError(err error) MyError {
	myError := MyError{
		Err: err,
	}
	if err != nil {
		myError.Message = Message{DefaultMessage: err.Error()}
	}
	return myError
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
