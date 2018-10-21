/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:26
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:22:12
 */
package tcplibrary

import "golang.org/x/net/websocket"

/* tcp库用到的接口定义 */

// Socket tcp通讯需要的一些回调函数
type Socket interface {
	OnConnect(*Conn) error           // 连接建立时
	OnError(error)                   // 连接发生错误
	OnClose(*Conn, error)            // 关闭连接时
	OnRecMessage(*Conn, interface{}) // 接收消息时
}

// ServerSocket 服务接口，实例化tcp server时传次参数
type ServerSocket interface {
	Socket
	GetClientID() string // 获取session id生成规则
}

// Packet 封包和解包
type Packet interface {
	Unmarshal(data []byte, c chan interface{}) ([]byte, error)                // 解包
	Marshal(v interface{}) ([]byte, error)                                    // 封包
	MarshalToJSON(v interface{}) (data []byte, payloadType byte, err error)   // 封包为json字符串形式
	UnmarshalToJSON(data []byte, payloadType byte, v interface{}) (err error) // 解包为json字符串形式
	GetWebsocketCodec() *websocket.Codec
}
