package types

import "github.com/google/gopacket"

// Sniff 嗅探任务
type Sniff interface {
	GetFilter() string        //  过滤器
	Callback(gopacket.Packet) //  回调
	Start() error             //  初始化任务
	End()                     // 销毁时任务
}
