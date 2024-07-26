package response

import "github.com/MrKrisYu/koi-go-common/sdk/i18n"

type Status struct {
	Code    int          `json:"code"`
	Message i18n.Message `json:"message"`
}

var (
	OK = Status{Code: 2000, Message: i18n.Message{ID: "response.ok", DefaultMessage: "OK"}}
)
