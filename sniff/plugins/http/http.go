package http

import (
	"bufio"
	"bytes"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/OhYee/rainbow/log"

	// mysql 数据库
	_ "github.com/go-sql-driver/mysql"
)

// Sniffer HTTP 嗅探器
type Sniffer struct {
	conn *sql.DB
}

// GetFilter 获取抓包过滤器
func (sniffer *Sniffer) GetFilter() string {
	return "dst port 80"
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

// Callback HTTP 回调函数
func (sniffer *Sniffer) Callback(pkg gopacket.Packet) {
	if pkg.TransportLayer() == nil || pkg.TransportLayer().LayerType() != layers.LayerTypeTCP {
		return
	}
	tcp := pkg.TransportLayer().(*layers.TCP)

	if len(tcp.Payload) < 4 ||
		!(strings.ToLower(string(tcp.Payload[:3])) == "get" ||
			strings.ToLower(string(tcp.Payload[:4])) == "post") {
		return
	}

	buf := bytes.NewBuffer(tcp.Payload)

	req, err := http.ReadRequest(bufio.NewReader(buf))
	if err != nil {
		return
	}
	now := time.Now()

	_, err = sniffer.conn.Exec("INSERT INTO `http` (`host`,`request`,`time`) values(?,?,?)", strings.ToLower(req.Host), string(tcp.Payload), now)
	if err != nil {
		log.Error.Println(err)
	}

}
