package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

const (
	// 字符集
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes = "0123456789"
	allBytes    = letterBytes + numberBytes

	// 用户名前缀
	usernamePrefix = "mall"

	// 随机部分长度
	randomLength = 8
)

// 使用 sync.Pool 来重用 strings.Builder
var builderPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	builder := builderPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		builderPool.Put(builder)
	}()

	builder.Grow(length)

	for i := 0; i < length; i++ {
		// 使用 crypto/rand 生成随机索引
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(allBytes))))
		if err != nil {
			// 如果出错，使用时间戳作为后备方案
			n = big.NewInt(time.Now().UnixNano() % int64(len(allBytes)))
		}
		builder.WriteByte(allBytes[n.Int64()])
	}

	return builder.String()
}

// GenerateUniqueUsername 生成唯一用户名
func GenerateUniqueUsername() string {
	builder := builderPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		builderPool.Put(builder)
	}()

	// 预分配足够的空间
	builder.Grow(len(usernamePrefix) + 10)

	// 写入前缀
	builder.WriteString(usernamePrefix)

	// 生成 6 位时间戳（微秒）
	timestamp := time.Now().UnixMicro() % 1000000
	builder.WriteString(fmt.Sprintf("%06d", timestamp))

	// 生成 4 位随机字符
	randomStr := GenerateRandomString(4)
	builder.WriteString(randomStr)

	return builder.String()
}

// GenerateRandomBytes 生成指定长度的随机字节
func GenerateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomInt 生成指定范围内的随机整数
func GenerateRandomInt(min, max int64) (int64, error) {
	if min >= max {
		return 0, fmt.Errorf("min must be less than max")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		return 0, err
	}

	return n.Int64() + min, nil
}

// GenerateRandomID 生成随机ID
func GenerateRandomID() string {
	// 使用 UUID v4
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// 如果出错，使用时间戳
		binary.BigEndian.PutUint64(b, uint64(time.Now().UnixNano()))
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
