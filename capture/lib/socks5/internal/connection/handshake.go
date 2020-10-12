package connection

import (
	"io"

	"github.com/OhYee/goutils/bytes"
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

// ClientHandshake Socks5 握手部分 — 客户端握手
//	版本      客户端支持的认证数       客户端支持的认证列表
//	1 字节    1 字节                 每个认证 1 字节
type ClientHandshake struct {
	version       byte
	verifyMethods []VerifyType
}

// ReadClientHandshake 读入客户端握手内容
func ReadClientHandshake(r io.Reader) (pkg ClientHandshake, err error) {
	b, err := bytes.ReadNBytes(r, 2)
	if err != nil {
		return
	}
	pkg = ClientHandshake{
		version:       b[0],
		verifyMethods: nil,
	}
	b, err = bytes.ReadNBytes(r, int(b[1]))
	if err != nil {
		return
	}
	pkg.verifyMethods = make([]VerifyType, len(b))
	for idx, v := range b {
		pkg.verifyMethods[idx] = VerifyType(v)
	}
	return
}

// ToBytes 将客户端握手转换为字节数组
func (pkg ClientHandshake) ToBytes() []byte {
	buf := bytes.NewBuffer(pkg.version, byte(len(pkg.verifyMethods)))
	for _, v := range pkg.verifyMethods {
		buf.WriteByte(byte(v))
	}
	return buf.Bytes()
}

// ServerHandshake Socks5 握手部分 - 服务端握手
// 版本     支持的认证方式
// 1 字节   1 字节
type ServerHandshake struct {
	version      byte
	verifyMethod VerifyType
}

// ReadServerHandshake 读入服务端握手内容
func ReadServerHandshake(r io.Reader) (pkg ServerHandshake, err error) {
	b, err := bytes.ReadNBytes(r, 2)
	if err != nil {
		return
	}
	pkg = ServerHandshake{
		version:      b[0],
		verifyMethod: VerifyType(b[1]),
	}
	return
}

// ToBytes 将服务端握手转换为字节数组
func (pkg ServerHandshake) ToBytes() []byte {
	return []byte{pkg.version, byte(pkg.verifyMethod)}
}
