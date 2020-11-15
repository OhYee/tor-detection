package command

import (
	"fmt"
	"io"

	"github.com/OhYee/goutils/bytes"
	"github.com/OhYee/rainbow/errors"
)

// AddressType 地址类型
// 0x01 IP V4地址
// 0x03 域名地址，域名地址的第1个字节为域名长度，剩下字节为域名名称字节数组
// 0x04 IP V6地址
type AddressType byte

const (
	// AddressIPv4 IPv4 地址
	AddressIPv4 AddressType = 1
	// AddressDomain 域名地址
	AddressDomain AddressType = 3
	// AddressIPv6 IPv6 地址
	AddressIPv6 AddressType = 4
)

var addressTypeName = [...]string{
	"Unknown",
	"AddressIPv4",
	"AddressDomain",
	"Unknown",
	"AddressIPv6",
}

// String 输出类型名称
func (at AddressType) String() string {
	if int(at) > len(addressTypeName) {
		return "Unknown"
	}
	return addressTypeName[at]
}

func (at AddressType) Read(r io.Reader) (b []byte, err error) {
	defer errors.Wrapper(&err)

	switch at {
	case AddressIPv4:
		b, err = bytes.ReadNBytes(r, 4)
		return
	case AddressDomain:
		b, err = bytes.ReadNBytes(r, 1)
		if err != nil {
			return
		}
		b, err = bytes.ReadNBytes(r, int(b[0]))
		return
	case AddressIPv6:
		b, err = bytes.ReadNBytes(r, 16)
		return
	default:
		err = errors.New("Can not parse AddressType %d", at)
		return
	}
}

// Address 将地址转换为可读形式
func (at AddressType) Address(address []byte) string {
	switch at {
	case AddressIPv4, AddressIPv6:
		return fmt.Sprintf("%v", address)
	default:
		return string(address)
	}
}
