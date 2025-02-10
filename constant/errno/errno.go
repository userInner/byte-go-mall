package errno

import "byte-go-mall/kitex_gen/base"

var (
	Success       = &base.Status{Code: 0, Message: "Success"}
	InternalError = &base.Status{Code: 0, Message: "InternalError"}
	DatabaseError = &base.Status{Code: 10002, Message: "Database error"}
	InvalidParams = &base.Status{Code: 10003, Message: "Invalid parameters"}

	UserExist     = &base.Status{Code: 20001, Message: "User already exists"}
	UserNotFound  = &base.Status{Code: 20002, Message: "User not found"}
	PasswordError = &base.Status{Code: 20003, Message: "Password error"}
	WeakPassword  = &base.Status{Code: 20004, Message: "Password is too weak"} // 用户密码太简单
)

// WithMessage 为Status添加自定义消息
func WithMessage(status *base.Status, msg string) *base.Status {
	return &base.Status{
		Code:    status.Code,
		Message: msg,
	}
}
