package command

import "github.com/OhYee/goutils/bytes"

// ResponseType 代理服务器响应类型
type ResponseType byte

const (
	// ResponseSuccess 代理服务器连接目标服务器成功
	ResponseSuccess ResponseType = iota
	// ResponseFailed 代理服务器故障
	ResponseFailed
	// ResponseRuleReject 代理服务器规则集不允许连接
	ResponseRuleReject
	// ResponseNetworkError 网络无法访问
	ResponseNetworkError
	// ResponseResolveFailed 主机名无效
	ResponseResolveFailed
	// ResponseReject 目标服务器拒绝
	ResponseReject
	// ResponseTTLExpire TTL 过期
	ResponseTTLExpire
	// ResponseUnknowCommand 不支持的命令
	ResponseUnknowCommand
	// ResponseUnknowaddress 不支持的目标服务器类型
	ResponseUnknowaddress
)

var responseTypeName = [...]string{
	"ResponseSuccess",
	"ResponseFailed",
	"ResponseRuleReject",
	"ResponseNetworkError",
	"ResponseResolveFailed",
	"ResponseReject",
	"ResponseTTLExpire",
	"ResponseUnknowCommand",
	"ResponseUnknowaddress",
}

// String 输出类型名称
func (rt ResponseType) String() string {
	if int(rt) > len(responseTypeName) {
		return "Unknown"
	}
	return responseTypeName[rt]
}

// CommandServer Socks5 代理服务器响应
// 版本	    响应	RSV	   地址类型	 地址	    端口号
// 1 字节	1 字节	1 字节	1 字节	1-255字节	2 字节
type CommandServer struct {
	Version     byte
	Response    ResponseType
	AddressType AddressType
	Address     []byte
	Port        int16
}

// Bytes 转换为二进制数组
func (cmd CommandServer) Bytes() []byte {
	buf := bytes.NewBuffer(cmd.Version, 0, byte(cmd.Response), 0, byte(cmd.AddressType))
	buf.Write(cmd.Address)
	buf.Write(bytes.FromInt16(cmd.Port))
	return buf.Bytes()
}
