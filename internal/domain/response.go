package domain

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string       `json:"error"`
}

func ErrorResp(msg string, err error) Response {
	resp := Response{
		Code:    1,
		Message: msg,
	}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp
}

func SuccessResp(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}
