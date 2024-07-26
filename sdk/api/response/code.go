package response

import "github.com/MrKrisYu/koi-go-common/sdk/i18n"

type Status struct {
	Code    int          `json:"code"`
	Message i18n.Message `json:"message"`
}

var (
	OK                  = Status{Code: 2000, Message: i18n.Message{ID: "response.ok", DefaultMessage: "OK"}}
	RequestParamsError  = Status{Code: 4000, Message: i18n.Message{ID: "response.badRequest", DefaultMessage: "Request Parameters Error"}}
	Unauthorized        = Status{Code: 4001, Message: i18n.Message{ID: "response.unauthorized", DefaultMessage: "Unauthorized Request"}}
	InternalServerError = Status{Code: 5000, Message: i18n.Message{ID: "response.internalServerError", DefaultMessage: "Internal Server Error"}}
	ServiceError        = Status{Code: 5001, Message: i18n.Message{ID: "response.serviceError", DefaultMessage: "Service Error"}}
)
