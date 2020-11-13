package dns

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/OhYee/rainbow/log"

	// mysql 数据库
	_ "github.com/go-sql-driver/mysql"
)

// Sniffer DNS 嗅探器
type Sniffer struct {
	conn *sql.DB
}

// GetFilter 获取抓包过滤器
func (sniffer *Sniffer) GetFilter() string {
	return "dst port 53"
}

// Start 程序开始运行时的任务
func (sniffer *Sniffer) Start() (err error) {
	sniffer.conn, err = sql.Open("mysql", "tor@tcp(127.0.0.1)/tor")
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

// Callback DNS 回调函数
func (sniffer *Sniffer) Callback(pkg gopacket.Packet) {
	if pkg.ApplicationLayer() == nil || pkg.ApplicationLayer().LayerType() != layers.LayerTypeDNS {
		return
	}
	dns := pkg.ApplicationLayer().(*layers.DNS)
	now := time.Now()
	for _, q := range dns.Questions {
		_, err := sniffer.conn.Exec("INSERT INTO `domain` (`domain`,`time`) values(?,?)", strings.ToLower(string(q.Name)), now)
		if err != nil {
			log.Error.Println(err)
		}
	}

}
