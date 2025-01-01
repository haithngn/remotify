package data

type AppErrorCode string

const (
	AppErrorCode_Invalid_Header   AppErrorCode = "1"
	AppErrorCode_Cannot_Make                   = "2"
	AppErrorCode_Invalid_Response              = "3"
)

var AppErrorCodes = []struct {
	Value  AppErrorCode
	TSName string
}{
	{AppErrorCode_Invalid_Header, "AppErrorCode_Invalid_Header"},
	{AppErrorCode_Cannot_Make, "AppErrorCode_Cannot_Make"},
	{AppErrorCode_Invalid_Response, "AppErrorCode_Invalid_Response"},
}
