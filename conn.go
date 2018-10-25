/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:02
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:09:56
 */
package tcplibrary

import (
	"net"
)

/* tcp 连接定义 */

// TCPType tcp连接类型
type TCPType = int

const (
	// TCPSocketType tcp连接
	TCPSocketType = iota
	// WebSocketType WebSocket连接
	WebSocketType
)

// Conn 自定义连接对象结构体,可以存储tcp或webSocket连接对象
type Conn struct {
	net.Conn
	connType TCPType // 连接对象类型
	clientID string  // 客户端id
	packet   Packet  // 封闭解包对象
}

// SendMessage 发送消息，参数为自己报文结构体
func (c *Conn) SendMessage(v interface{}) (int, error) {
	// 判断是tcp还是websocket
	if c.connType == TCPSocketType { // 二进制协议
		// 先封包，再发送数据
		data, err := c.packet.Marshal(v)
		if err != nil {
			globalLogger.Errorf(err.Error())
			return 0, err
		}
		return c.Write(data)
	} else if c.connType == WebSocketType { // json方式
		data, _, err := c.packet.MarshalToJSON(v)
		if err != nil {
			globalLogger.Errorf(err.Error())
			return 0, err
		}
		c.Write(data)
	} else {
		globalLogger.Errorf("不支持的连接方式")
	}
	return 0, nil
}

// GetClientID 获取当前连接id
func (c *Conn) GetClientID() string {
	return c.clientID
}

// GetConnType 获取连接类型
func (c *Conn) GetConnType() TCPType {
	return c.connType
}
