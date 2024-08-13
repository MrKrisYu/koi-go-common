package i18n

type MyError struct {
	Err      error
	Messages []Message `json:"messages"`
}

func (e *MyError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "<nil>"
}

func (e *MyError) AddMessage(mi MessageID, arg ...any) {
	if e.Messages == nil {
		e.Messages = []Message{}
	}
	newMsg := NewMessage(mi, arg...)
	e.Messages = append(e.Messages, newMsg)
}

func NewMyError(err error) *MyError {
	myError := MyError{
		Err:      err,
		Messages: []Message{},
	}
	if err != nil {
		myError.Messages = append(myError.Messages, Message{DefaultMessage: err.Error()})
	}
	return &myError
}

func NewMyErrorWithMessageID(mi MessageID, arg ...any) *MyError {
	myError := MyError{
		Err:      nil,
		Messages: []Message{},
	}
	newMsg := NewMessage(mi, arg...)
	myError.Messages = append(myError.Messages, newMsg)
	return &myError
}
