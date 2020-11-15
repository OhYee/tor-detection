package server

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/OhYee/tor-detection/capture/socks5/command"
	"github.com/OhYee/tor-detection/capture/socks5/handshake"

	"github.com/OhYee/rainbow/errors"
	"github.com/OhYee/rainbow/log"
)

type CaptureHandle func(id string, src string, dst string) (rw io.ReadWriter, localAddress string, err error)

func HandleConn(conn net.Conn, captureHandle CaptureHandle) {
	c := NewConnection(conn, captureHandle)
	err := c.Serve()
	if err != nil {
		log.Error.Println(errors.ShowStack(err))
	}
}

type Connection struct {
	conn          net.Conn
	signal        chan bool
	remoteConn    io.ReadWriter
	captureHandle CaptureHandle
}

func NewConnection(conn net.Conn, captureHandle CaptureHandle) *Connection {
	return &Connection{conn: conn, captureHandle: captureHandle}
}

func (conn *Connection) GetConnection() net.Conn {
	return conn.conn
}

func (conn *Connection) Serve() error {
	c := conn.GetConnection()
	/*
		握手部分
	*/
	clientHandshake, err := handshake.ReadClientHandshake(c)
	if err != nil {
		return err
	}
	if clientHandshake.Version != 5 {
		return errors.New("Can only support Socks5 now")
	}

	verifyMethod := handshake.ChooseVerify(clientHandshake.VerifyMethods)
	serverHandshake := handshake.ServerHandshake{
		Version:      clientHandshake.Version,
		VerifyMethod: verifyMethod,
	}
	c.Write(serverHandshake.ToBytes())

	err = verifyMethod.Verify(c)
	if err != nil {
		return err
	}

	cmd, err := command.ReadCommand(c)
	if err != nil {
		return err
	}

	if cmd.Version != 5 {
		return errors.New("Can only support Socks5 now")
	}

	log.Debug.Printf("Socks5 [%s] to %s:%d", cmd.Command.String(), cmd.AddressType.Address(cmd.Address), cmd.Port)

	/*
		Socks5 代理服务器响应
		版本	 响应	RSV	   地址类型	 地址	    端口号
		1 字节	1 字节	1 字节	1 字节	1-255字节	2 字节

		0x00 代理服务器连接目标服务器成功
		0x01 代理服务器故障
		0x02 代理服务器规则集不允许连接
		0x03 网络无法访问
		0x04 目标服务器无法访问（主机名无效）
		0x05 连接目标服务器被拒绝
		0x06 TTL已过期
		0x07 不支持的命令
		0x08 不支持的目标服务器地址类型
		0x09 - 0xFF 未分配
	*/
	var localAddress string
	conn.remoteConn, localAddress, err = conn.captureHandle(fmt.Sprintf("%+v", &conn), conn.conn.RemoteAddr().String(), conn.conn.LocalAddr().String())
	if err != nil {
		return errors.NewErr(err)
	}

	localIP, localPort, err := net.SplitHostPort(localAddress)
	if err != nil {
		return errors.NewErr(err)
	}

	localPortInt, err := strconv.Atoi(localPort)
	if err != nil {
		return errors.NewErr(err)
	}

	commandResponse := command.CommandServer{
		Version:     5,
		Response:    command.ResponseSuccess,
		AddressType: command.AddressIPv4,
		Address:     net.ParseIP(localIP)[:4],
		Port:        int16(localPortInt),
	}
	c.Write(commandResponse.Bytes())

	go conn.redirect(conn.conn, conn.remoteConn, "Client -%d-> Socks5 -%d-> Server\n")
	go conn.redirect(conn.remoteConn, conn.conn, "Client <-%d- Socks5 <-%d- Server\n")

	<-conn.signal
	return nil
}

func (conn Connection) redirect(src io.ReadWriter, dst io.ReadWriter, logFormat string) {
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if err != nil {
			if io.EOF != err {
				log.Error.Println(errors.ShowStack(errors.NewErr(err)))
			} else {
				dst.Write(buf[:n])
			}
			conn.Close()
			return
		}

		nn, err := dst.Write(buf[:n])
		if err != nil {
			if io.EOF != err {
				log.Error.Println(errors.ShowStack(errors.NewErr(err)))
			}
			conn.Close()
			return
		}

		log.Debug.Printf("%+v", time.Now().UnixNano())
		log.Debug.Printf(logFormat, n, nn)
	}
}

func (conn Connection) Close() {
	conn.signal <- true
}
