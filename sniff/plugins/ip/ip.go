package ip

import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/OhYee/rainbow/log"

	// mysql 数据库
	_ "github.com/go-sql-driver/mysql"
)

// Sniffer  嗅探器
type Sniffer struct {
	conn *sql.DB
	ips  []net.IP
}

// GetFilter 获取抓包过滤器
func (sniffer *Sniffer) GetFilter() string {
	return strings.Join(getIPs(), " or ")
}

// Start 程序开始运行时的任务
func (sniffer *Sniffer) Start() (err error) {
	sniffer.conn, err = sql.Open("mysql", "tor@tcp(127.0.0.1)/tor")
	ips := getIPs()
	sniffer.ips = make([]net.IP, len(ips))
	for idx, ip := range ips {
		sniffer.ips[idx] = net.ParseIP(ip)
	}
	return
}

// End 关闭数据库连接
func (sniffer *Sniffer) End() {
	if sniffer.conn != nil {
		if err := sniffer.conn.Close(); err != nil {
			log.Error.Println(err)
		}
	}
}

// Callback  回调函数
func (sniffer *Sniffer) Callback(pkg gopacket.Packet) {
	if pkg.NetworkLayer() == nil {
		return
	}

	var dst net.IP = nil
	if pkg.NetworkLayer().LayerType() == layers.LayerTypeIPv4 {
		ipv4 := pkg.NetworkLayer().(*layers.IPv4)
		dst = ipv4.DstIP

	}
	if pkg.NetworkLayer().LayerType() == layers.LayerTypeIPv6 {
		ipv6 := pkg.NetworkLayer().(*layers.IPv6)
		dst = ipv6.DstIP
	}
	if dst != nil && sniffer.checkIP(dst) {
		_, err := sniffer.conn.Exec("INSERT INTO `ip` (`ip`, `count`) VALUES (?, 1) ON DUPLICATE KEY UPDATE `count` = `count` + 1", dst.String())
		if err != nil {
			log.Error.Println(err)
		}
	}
}

func (sniffer *Sniffer) checkIP(ip net.IP) bool {
	for _, t := range sniffer.ips {
		if t.Equal(ip) {
			return false
		}
	}
	return true
}

func getIPs() []string {
	addresses, _ := net.InterfaceAddrs()
	ips := make([]string, 0)
	for _, ip := range addresses {
		realIP := strings.Split(ip.String(), "/")[0]
		if realIP != "127.0.0.1" && realIP != "::1" {
			ips = append(ips, fmt.Sprintf("src %s", realIP))
		}
	}
	return ips
}
