package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"

	gnet "github.com/OhYee/goutils/net"
	"github.com/OhYee/rainbow/errors"
	"github.com/OhYee/tor-detection/capture/lib/log"
	"github.com/OhYee/tor-detection/capture/lib/socks5"
)

var (
	file    *os.File = os.Stdout
	payload          = false
)

type Capture struct {
	id         string
	srcIP      net.IP
	srcPort    uint16
	dstIP      net.IP
	dstPort    uint16
	remoteConn net.Conn
}

func newCapture(id string, src string, dst string) (capture io.ReadWriter, localAddress string, err error) {
	srcIP, srcPort, err := gnet.ParseIPv4Port(src)
	if err != nil {
		return
	}
	dstIP, dstPort, err := gnet.ParseIPv4Port(dst)
	if err != nil {
		return
	}
	remoteConn, err := net.Dial("tcp", dst)
	if err != nil {
		return
	}
	localAddress = remoteConn.LocalAddr().String()
	capture = Capture{
		id:         id,
		srcIP:      srcIP,
		srcPort:    srcPort,
		dstIP:      dstIP,
		dstPort:    dstPort,
		remoteConn: remoteConn,
	}
	return
}
func (cap Capture) Read(p []byte) (n int, err error) {
	n, err = cap.remoteConn.Read(p)
	if err != nil {
		return
	}
	file.Write([]byte(fmt.Sprintf("\nr %s %d\n", time.Now().Format("2006-01-02 15:04:05"), n)))
	if payload {
		file.Write(p[:n])
	}
	return
}
func (cap Capture) Write(p []byte) (n int, err error) {
	n, err = cap.remoteConn.Write(p)
	if err != nil {
		return
	}
	file.Write([]byte(fmt.Sprintf("\nw %s %d\n", time.Now().Format("2006-01-02 15:04:05"), n)))
	if payload {
		file.Write(p[:n])
	}
	return
}

func main() {
	var debug, silent bool
	var ip, output string
	var port int
	var err error

	flag.BoolVar(&debug, "d", false, "Debug mode")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.BoolVar(&silent, "s", false, "Silent mode(Only show error)")
	flag.BoolVar(&silent, "silent", false, "Silent mode(Only show error)")

	flag.StringVar(&ip, "ip", "127.0.0.1", "Listen ip address")
	flag.IntVar(&port, "p", 1080, "Listen port")
	flag.IntVar(&port, "port", 1080, "Listen port")

	flag.BoolVar(&payload, "payload", false, "Record payload")
	flag.StringVar(&output, "o", "", "Record file position")
	flag.StringVar(&output, "output", "", "Record file position")

	flag.Parse()

	if silent {
		log.Error.SetOutputToStdout()
		log.Info.SetOutputToNil()
		log.Debug.SetOutputToNil()
	} else {
		if debug {
			log.Error.SetOutputToStdout()
			log.Info.SetOutputToStdout()
			log.Debug.SetOutputToStdout()
		} else {
			log.Error.SetOutputToStdout()
			log.Info.SetOutputToStdout()
			log.Debug.SetOutputToNil()
		}
	}

	if output != "" {
		file, err = os.Create(output)
		if err != nil {
			panic(err)
		}
	}

	socks := socks5.NewSocks5Server(newCapture)
	err = socks.Start(ip, port)
	if err != nil {
		log.Error.Println(errors.ShowStack(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	file.Close()
}
