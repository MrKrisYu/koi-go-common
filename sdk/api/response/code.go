package response

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	OK                    = Status{Code: 2000, Message: "OK"}
	RequestParamsError    = Status{Code: 4000, Message: "Request Parameters Error"}
	Unauthorized          = Status{Code: 4001, Message: "Unauthorized Request"}
	NotLoginUserOperation = Status{Code: 4002, Message: "Error for overstepping your authority"}
	InternalServerError   = Status{Code: 5000, Message: "Internal Server Error"}
	ServiceError          = Status{Code: 5001, Message: "Service Error"}
)
