package main

//go:generate go build -o ../../bin/mark mark.go
import (
	"database/sql"
	"flag"
	"fmt"
	"net"
	"net/textproto"
	"strings"

	"github.com/cretz/bine/control"
	// mysql 数据库
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	password := ""
	port := 9051
	flag.StringVar(&password, "password", "", "password of tor controller")
	flag.IntVar(&port, "port", 9051, "port of tor controller")
	flag.IntVar(&port, "p", 9051, "port of tor controller")
	flag.Parse()

	tcpConn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	defer tcpConn.Close()
	conn := control.NewConn(textproto.NewConn(tcpConn))
	defer conn.Close()
	conn.Authenticate(password)
	res, err := conn.GetInfo("ns/all")
	if err != nil {
		panic(err)
	}
	ips := make([]string, 0)
	for _, r := range res {
		if r.Key == "ns/all" {
			lines := strings.Split(r.Val, "\n")
			for _, line := range lines {
				part := strings.Split(line, " ")
				if part[0] == "r" {
					ips = append(ips, part[6])
				}
			}
		}
	}
	db, err := sql.Open("mysql", "tor@tcp(127.0.0.1)/tor")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	l := len(ips)
	for idx, ip := range ips {
		fmt.Printf("%d/%d\r", idx+1, l)
		_, err = db.Exec("UPDATE `tor`.`ip` SET `tor`=1 WHERE `ip`=?", ip)
		if err != nil {
			fmt.Println(err)
		}
	}
}
