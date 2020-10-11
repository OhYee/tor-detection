package socks5

import (
	"net"

	"github.com/OhYee/tor-detection/capture/lib/log"
	"github.com/OhYee/tor-detection/capture/lib/socks5/internal/connection"
)

type Socks5Server struct {
	captureHandle CaptureHandle
	signal        chan bool
}
type CaptureHandle = connection.CaptureHandle

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
				go connection.HandleConn(conn, socks.captureHandle)
			}
		}

	}
}
