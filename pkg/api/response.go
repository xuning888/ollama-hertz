package api

type ResponseData struct {
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	RedirectUrl string      `json:"redirectUrl"`
	Data        interface{} `json:"data"`
}

var (
	SuccessResponseCode = &ResponseCode{Code: "S00000", Message: "success"}
	SystemError         = &ResponseCode{Code: "E00001", Message: "system error"}
	ParamsError         = &ResponseData{Code: "E00002", Message: "params error"}
)

type ResponseCode struct {
	Code    string
	Message string
}

func WithResponseCode(responseCode *ResponseCode) *ResponseData {
	return &ResponseData{
		Code:    responseCode.Code,
		Message: responseCode.Message,
	}
}

func SuccessWithData(data interface{}) *ResponseData {
	responseData := Success()
	responseData.Data = data
	return responseData
}

func Success() *ResponseData {
	return WithResponseCode(SuccessResponseCode)
}

func Failed() *ResponseData {
	return WithResponseCode(SystemError)
}

func FailedWithMessage(message string) *ResponseData {
	errorResp := WithResponseCode(SystemError)
	errorResp.Message = message
	return errorResp
}
