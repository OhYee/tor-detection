package command

// CommandType 命令类型：
// 0x01 CONNECT 连接上游服务器
// 0x02 BIND 绑定，客户端会接收来自代理服务器的链接，著名的FTP被动模式
// 0x03 UDP ASSOCIATE UDP中继
type CommandType byte

const (
	// CommandConnect 连接上游服务器
	CommandConnect CommandType = 1
	// CommandBind 绑定（接收来自服务器的连接）
	CommandBind CommandType = 2
	// CommandUDP UDP 中继
	CommandUDP CommandType = 3
)

var commandTypeName = [...]string{
	"Unknown",
	"CommandConnect",
	"CommandBind",
	"CommandUDP",
}

// String 输出类型名称
func (ct CommandType) String() string {
	if int(ct) > len(commandTypeName) {
		return "Unknown"
	}
	return commandTypeName[ct]
}
