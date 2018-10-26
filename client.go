/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:54:57
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:11:34
 */
package tcplibrary

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"time"
)

/* tcp golang客户端 */

// TCPClient tcp客户端
type TCPClient struct {
	*TCPLibrary
	pingData interface{} // ping时的包
	conn     *Conn       // 连接对象
}

// NewTCPClient 创建一个tcp客户端
func NewTCPClient(debug bool, socket Socket, packets ...Packet) (*TCPClient, error) {
	if socket == nil {
		return nil, errors.New("Socket参数不能是nil")
	}
	// 封包解包对象
	var packet Packet
	if len(packets) == 0 {
		packet = NewDefaultPacket()
	} else {
		packet = packets[0]
	}

	// 请求上下文
	ctx, cancel := context.WithCancel(context.Background())

	return &TCPClient{
		TCPLibrary: &TCPLibrary{
			ctx:            ctx,
			cancel:         cancel,
			packet:         packet,
			socket:         socket,
			readDeadline:   DefaultReadDeadline,
			readBufferSize: DefaultBufferSize,
			isServer:       false,
		},
	}, nil
}

// DialAndStart 连接到服务器，并开始读取信息
func (c *TCPClient) DialAndStart(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		globalLogger.Errorf(err.Error())
		return err
	}

	return c.Dial(conn)
}

// DialAndStartTLS 连接到服务器，并开始读取信息 tls
func (c *TCPClient) DialAndStartTLS(address string) error {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", address, conf)
	if err != nil {
		globalLogger.Errorf(err.Error())
		return err
	}

	return c.Dial(conn)
}

// Dial 连接tcp服务端
func (c *TCPClient) Dial(conn net.Conn) error {
	// 判断是否设置读超时
	if c.readDeadline == 0 {
		c.readDeadline = DefaultReadDeadline
	}
	// 赋值给当前连接对象
	c.conn = &Conn{
		Conn:     conn,
		connType: TCPSocketType,
		packet:   c.packet,
	}
	// 通知建立连接
	err := c.socket.OnConnect(c.conn)
	if err != nil {
		globalLogger.Errorf(err.Error())
		// 如果建立连接函数返回false，则关闭连接
		c.socket.OnClose(c.conn, err) // 通知关闭
		err = conn.Close()            // 关闭连接
		if err != nil {
			globalLogger.Errorf(err.Error())
		}
		return err
	}
	// 开启一个协程处理数据接收
	go c.handleConn(c.conn)
	return nil
}

// GetConn 获取连接对象
func (c *TCPClient) GetConn() *Conn {
	return c.conn
}

// Ping 保活
func (c *TCPClient) Ping(v interface{}) {
	go func() {
		for {
			_, err := c.conn.SendMessage(v)
			if err != nil {
				globalLogger.Errorf("client ping error", err)
				return
			}
			if DefaultReadDeadline < 3*time.Second {
				time.Sleep(1 * time.Second)
			} else {
				time.Sleep(DefaultReadDeadline / 3 * 2)
			}
		}
	}()
}
