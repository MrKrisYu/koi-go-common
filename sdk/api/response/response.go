package response

type Response[T any] struct {
	Code    int    `json:"code"`    // 状态响应码
	Message string `json:"message"` // 响应消息
	Data    T      `json:"data"`    // 响应数据
}
