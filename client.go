/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:54:57
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:11:34
 */
package tcplibrary

import (
	"crypto/tls"
	"errors"
	"net"
)

/* tcp golang客户端 */

// TCPClient tcp客户端
type TCPClient struct {
	*TCPLibrary
	conn *Conn // 连接对象
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
	// 标记为客户端
	isServer = false

	return &TCPClient{
		TCPLibrary: &TCPLibrary{
			packet:         packet,
			socket:         socket,
			readDeadline:   DefaultReadDeadline,
			readBufferSize: DefaultBufferSize,
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

// DialAndStart 连接到服务器，并开始读取信息
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
