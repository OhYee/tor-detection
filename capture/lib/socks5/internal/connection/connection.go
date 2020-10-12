package connection

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/OhYee/rainbow/errors"
	"github.com/OhYee/tor-detection/capture/lib/log"

	"github.com/OhYee/goutils/bytes"
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
	clientHandshake, err := ReadClientHandshake(c)
	if err != nil {
		return err
	}
	if clientHandshake.version != 5 {
		return errors.New("Can only support Socks5 now")
	}

	verifyMethod := ChooseVerify(clientHandshake.verifyMethods)
	serverHandshake := ServerHandshake{
		version:      clientHandshake.version,
		verifyMethod: verifyMethod,
	}
	c.Write(serverHandshake.ToBytes())

	err = verifyMethod.Verify(c)
	if err != nil {
		return err
	}

	/*
		Socks5 命令
		版本号  命令    保留字段  目标服务器地址类型  目标服务器   端口号
		1 字节  1 字节  1 字节   1 字节            1~255 字节  2 字节

		命令类型：
		0x01 CONNECT 连接上游服务器
		0x02 BIND 绑定，客户端会接收来自代理服务器的链接，著名的FTP被动模式
		0x03 UDP ASSOCIATE UDP中继

		服务器类型：
		0x01 IP V4地址
		0x03 域名地址，域名地址的第1个字节为域名长度，剩下字节为域名名称字节数组
		0x04 IP V6地址
	*/

	b, err := bytes.ReadNBytes(c, 4)
	if err != nil {
		return err
	}

	version := b[0]
	if version != 5 {
		return errors.New("Can only support Socks5 now, got %d", version)
	}
	command := b[1]
	if command != 1 {
		return errors.New("Can only support CONNECT command")
	}
	addressType := b[3]

	address := ""
	switch addressType {
	case 1:
		b, err = bytes.ReadNBytes(c, 4)
		if err != nil {
			return errors.NewErr(err)
		}
		address = fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
	case 3:
		b, err = bytes.ReadNBytes(c, 1)
		if err != nil {
			return errors.NewErr(err)
		}
		length := b[0]
		b, err = bytes.ReadNBytes(c, int(length))
		if err != nil {
			return errors.NewErr(err)
		}
		ip, err := net.LookupIP(string(b))
		if err != nil {
			return errors.NewErr(err)
		}
		if len(ip) == 0 {
			return errors.New("Can not lookup the domain %s", string(b))
		}
		address = ip[0].String()
	case 4:
		return errors.New("Do not support IPv6")
	}

	b, err = bytes.ReadNBytes(c, 2)
	if err != nil {
		return errors.NewErr(err)
	}
	port := bytes.ToInt16(b)
	log.Debug.Printf("Socks5 to %+v:%d\n", address, port)

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
	conn.remoteConn, localAddress, err = conn.captureHandle(fmt.Sprintf("%+v", &conn), conn.conn.RemoteAddr().String(), fmt.Sprintf("%s:%d", address, port))
	// net.DialTCP("tcp", nil, &net.TCPAddr{
	// 	IP:   net.ParseIP(address),
	// 	Port: int(port),
	// })
	if err != nil {
		return errors.NewErr(err)
	}

	localIP, localPort, err := net.SplitHostPort(localAddress)
	if err != nil {
		return errors.NewErr(err)
	}

	c.Write([]byte{5, 0, 0, 1})
	c.Write([]byte(net.ParseIP(localIP))[:4])
	localPortInt, err := strconv.Atoi(localPort)
	if err != nil {
		return errors.NewErr(err)
	}
	c.Write(bytes.FromInt16(int16(localPortInt)))

	go conn.redirect(conn.conn, conn.remoteConn, "Client -%d-> Socks5 -%d-> Server\n")
	go conn.redirect(conn.remoteConn, conn.conn, "Client <-%d- Socks5 <-%d- Server\n")

	<-conn.signal
	return nil
}

// Handshake Socks5 握手阶段
func (conn *Connection) handshake() error {
	c := conn.GetConnection()
	b, err := bytes.ReadNBytes(c, 2)
	if err != nil {
		return err
	}

	/*
		Socks5 握手部分 — 客户端握手
		版本      客户端支持的认证数       客户端支持的认证列表
		1 字节    1 字节                 每个认证 1 字节
	*/
	version := b[0]
	if version != 5 {
		return errors.New("Can only support Socks5 now")
	}
	// 认证方式
	// 0x00 不需要认证（常用）
	// 0x01 GSSAPI认证
	// 0x02 账号密码认证（常用）
	// 0x03 - 0x7F IANA分配
	// 0x80 - 0xFE 私有方法保留
	// 0xFF 无支持的认证方法
	methodsCount := b[1]
	b, err = bytes.ReadNBytes(c, int(methodsCount))
	if err != nil {
		return errors.NewErr(err)
	}

	noVerify := false
	for _, method := range b {
		if method == 0 {
			noVerify = true
		}
	}
	if !noVerify {
		return errors.New("Client do not support no verify method")
	}

	/*
		Socks5 握手部分 - 服务端握手
		版本     支持的认证方式
		1 字节   1 字节
	*/
	c.Write([]byte{5, 0})

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
