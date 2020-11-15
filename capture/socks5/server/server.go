package server

import (
	"net"

	"github.com/OhYee/rainbow/log"
)

type Socks5Server struct {
	captureHandle CaptureHandle
	signal        chan bool
}

func NewSocks5Server(captureHandle CaptureHandle) *Socks5Server {
	return &Socks5Server{
		captureHandle: captureHandle,
		signal:        make(chan bool),
	}
}

func (socks *Socks5Server) Start(ip string, port int) (err error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		return err
	}
	log.Info.Printf("Socks5 server started at %s:%d\n", ip, port)
	for {
		select {
		case signal := <-socks.signal:
			if signal == true {
				break
			}
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Error.Println(err)
				continue
			} else {
				log.Info.Printf("New connection from %s\n", conn.RemoteAddr().String())
				go HandleConn(conn, socks.captureHandle)
			}
		}

	}
}
