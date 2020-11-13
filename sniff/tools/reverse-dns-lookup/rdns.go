package main

//go:generate go build -o ../../bin/rdns rdns.go
import (
	"database/sql"
	"fmt"
	"net"
	"strings"

	// mysql 数据库
	_ "github.com/go-sql-driver/mysql"
)

func reverseDNSLookup(ip string) string {
	domains, err := net.LookupAddr(ip)
	if err != nil {
		return "unknown"
	}
	return strings.Join(domains, ",")
}

func main() {
	conn, err := sql.Open("mysql", "tor@tcp(127.0.0.1)/tor")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT `ip` FROM `tor`.`ip` WHERE `domain` IS NULL")
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	ips := make([]string, 0)
	for rows.Next() {
		var ip string
		rows.Scan(&ip)
		ips = append(ips, ip)
	}

	l := len(ips)
	for idx, ip := range ips {
		fmt.Printf("%d/%d\r", idx, l)
		conn.Exec("UPDATE `tor`.`ip` SET `domain`=? where `ip`=?", reverseDNSLookup(ip), ip)
	}
}
