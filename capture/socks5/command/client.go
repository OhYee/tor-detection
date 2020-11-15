package command

import (
	"io"

	"github.com/OhYee/goutils/bytes"
	"github.com/OhYee/rainbow/errors"
)

// CommandClient Socks5 命令包
// 版本号  命令    保留字段  目标服务器地址类型  目标服务器   端口号
// 1 字节  1 字节  1 字节   1 字节            1~255 字节  2 字节
type CommandClient struct {
	Version     byte
	Command     CommandType
	AddressType AddressType
	Address     []byte
	Port        int16
}

// ReadCommand 读入 Socks5 命令包
func ReadCommand(r io.Reader) (cmd CommandClient, err error) {
	defer errors.Wrapper(&err)

	b, err := bytes.ReadNBytes(r, 4)
	cmd = CommandClient{
		Version:     b[0],
		Command:     CommandType(b[1]),
		AddressType: AddressType(b[3]),
		Address:     nil,
		Port:        0,
	}
	cmd.Address, err = cmd.AddressType.Read(r)
	if err != nil {
		return
	}

	b, err = bytes.ReadNBytes(r, 2)
	if err != nil {
		return
	}
	cmd.Port = bytes.ToInt16(b)
	return
}

// Bytes 转换为二进制数组
func (cmd CommandClient) Bytes() []byte {
	buf := bytes.NewBuffer(cmd.Version, byte(cmd.Command), 0, byte(cmd.AddressType))
	buf.Write(cmd.Address)
	buf.Write(bytes.FromInt16(cmd.Port))
	return buf.Bytes()
}
