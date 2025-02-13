package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	// 邮箱正则表达式
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// 手机号正则表达式（中国大陆）
	phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

	// 用户名正则表达式（字母、数字、下划线，长度3-32）
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

	// 基本密码格式（字母、数字和特殊字符，长度6-32）
	passwordRegex = regexp.MustCompile(`^[a-zA-Z\d!@#$%^&*]{6,32}$`)
)

// IsValidPassword 验证密码强度
func IsValidPassword(password string) bool {
	if password == "" {
		return false
	}

	// 长度检查
	if len(password) < 6 || len(password) > 32 {
		return false
	}

	// 基本格式检查
	if !passwordRegex.MatchString(password) {
		return false
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}

		// 必须包含大写字母、小写字母和数字
	}
	return hasUpper && hasLower && hasNumber
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}

	// 转换为小写
	email = strings.ToLower(email)

	// 长度检查
	if len(email) > 254 {
		return false
	}

	// 正则匹配
	return emailRegex.MatchString(email)
}

// IsValidPhone 验证手机号格式
func IsValidPhone(phone string) bool {
	if phone == "" {
		return false
	}

	// 去除空格和特殊字符
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, " ", "")

	// 正则匹配
	return phoneRegex.MatchString(phone)
}
