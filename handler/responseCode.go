package handler

type ResponseCode int

const (
	SuccessCode ResponseCode = 1000 + iota
	FailedCode
)

var codeMsgMap = map[ResponseCode]string{
	SuccessCode: "success",
	FailedCode:  "fail to insert",
}
