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

type StatusError struct {
	Status *base.Status
}

// Error 实现 error 接口的 Error 方法
func (se *StatusError) Error() string {
	return se.Status.Message
}

// WithMessage 为Status添加自定义消息
func WithMessage(status *base.Status, msg string) *base.Status {
	return &base.Status{
		Code:    status.Code,
		Message: msg,
	}
}

func NewError(status *base.Status, msg string) error {
	return &StatusError{
		Status: &base.Status{
			Code:    status.Code,
			Message: msg,
		},
	}
}
