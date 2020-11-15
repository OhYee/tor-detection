package handshake

import (
	"io"

	"github.com/OhYee/goutils/bytes"
	"github.com/OhYee/rainbow/errors"
)

// ServerHandshake Socks5 握手部分 - 服务端握手
// 版本     支持的认证方式
// 1 字节   1 字节
type ServerHandshake struct {
	Version      byte
	VerifyMethod VerifyType
}

// ReadServerHandshake 读入服务端握手内容
func ReadServerHandshake(r io.Reader) (pkg ServerHandshake, err error) {
	defer errors.Wrapper(&err)

	b, err := bytes.ReadNBytes(r, 2)
	if err != nil {
		return
	}
	pkg = ServerHandshake{
		Version:      b[0],
		VerifyMethod: VerifyType(b[1]),
	}
	return
}

// ToBytes 将服务端握手转换为字节数组
func (pkg ServerHandshake) ToBytes() []byte {
	return []byte{pkg.Version, byte(pkg.VerifyMethod)}
}
