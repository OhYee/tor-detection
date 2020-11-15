package handshake

import (
	"io"

	"github.com/OhYee/rainbow/errors"
)

// VerifyType 认证方式
// 0x00 不需要认证（常用）
// 0x01 GSSAPI认证
// 0x02 账号密码认证（常用）
// 0x03 - 0x7F IANA分配
// 0x80 - 0xFE 私有方法保留
// 0xFF 无支持的认证方法
type VerifyType byte

const (
	// VerifyNone 不需要认证
	VerifyNone VerifyType = 0
	// VerifyGSSAPI 使用 GSSAPI 认账
	VerifyGSSAPI VerifyType = 1
	// VerifyPassword 使用用户密码认证
	VerifyPassword VerifyType = 2
)

var typeName = [...]string{
	"VerifyNone",
	"VerifyGSSAPI",
	"VerifyPassword",
}

// String 输出类型名称
func (v VerifyType) String() string {
	if int(v) > len(typeName) {
		return "Unknown"
	}
	return typeName[v]
}

// Verify 验证用户身份
func (v VerifyType) Verify(r io.Reader) (err error) {
	defer errors.Wrapper(&err)

	switch v {
	case VerifyNone:
		return
	default:
		err = errors.New("Unknown verify method %d", v)
		return
	}
}

// ChooseVerify 从给定列表中选择支持的认证方式
func ChooseVerify(v []VerifyType) VerifyType {
	for _, t := range v {
		switch t {
		case VerifyNone:
			return t
		}
	}
	return VerifyNone
}
