package handshake

import (
	"io"

	"github.com/OhYee/goutils/bytes"
	"github.com/OhYee/rainbow/errors"
)

// ClientHandshake Socks5 握手部分 — 客户端握手
//	版本      客户端支持的认证数       客户端支持的认证列表
//	1 字节    1 字节                 每个认证 1 字节
type ClientHandshake struct {
	Version       byte
	VerifyMethods []VerifyType
}

// ReadClientHandshake 读入客户端握手内容
func ReadClientHandshake(r io.Reader) (pkg ClientHandshake, err error) {
	defer errors.Wrapper(&err)

	b, err := bytes.ReadNBytes(r, 2)
	if err != nil {
		return
	}
	pkg = ClientHandshake{
		Version:       b[0],
		VerifyMethods: nil,
	}
	b, err = bytes.ReadNBytes(r, int(b[1]))
	if err != nil {
		return
	}
	pkg.VerifyMethods = make([]VerifyType, len(b))
	for idx, v := range b {
		pkg.VerifyMethods[idx] = VerifyType(v)
	}
	return
}

// ToBytes 将客户端握手转换为字节数组
func (pkg ClientHandshake) ToBytes() []byte {
	buf := bytes.NewBuffer(pkg.Version, byte(len(pkg.VerifyMethods)))
	for _, v := range pkg.VerifyMethods {
		buf.WriteByte(byte(v))
	}
	return buf.Bytes()
}
