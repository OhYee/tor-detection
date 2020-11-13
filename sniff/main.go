package main

import (
	"fmt"
	"strings"

	"github.com/OhYee/tor-detection/sniff/plugins/dns"
	"github.com/OhYee/tor-detection/sniff/plugins/http"
	"github.com/OhYee/tor-detection/sniff/plugins/ip"

	"github.com/OhYee/tor-detection/sniff/lib/types"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var sniffs = []types.Sniff{
	&dns.Sniffer{},
	&http.Sniffer{},
	&ip.Sniffer{},
}

func main() {
	filters := make([]string, len(sniffs))
	for idx, s := range sniffs {
		filters[idx] = fmt.Sprintf("(%s)", s.GetFilter())
		if err := s.Start(); err != nil {
			panic(err)
		}
		defer s.End()
	}
	filter := strings.Join(filters, "or")
	fmt.Printf("filter: %s\n", filter)

	//打开网络接口，抓取在线数据
	handle, err := pcap.OpenLive("eth0", int32(65535), true, pcap.BlockForever)
	if err != nil {
		panic(fmt.Sprintln("启动失败", err))
	}
	// 设置过滤器
	if err := handle.SetBPFFilter(filter); err != nil {
		panic(fmt.Sprintln("过滤器设置失败", err))
	}
	defer handle.Close()

	fmt.Println("抓包开始")

	// 抓包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true
	for packet := range packetSource.Packets() {
		for _, s := range sniffs {
			s.Callback(packet)
		}
	}
}
