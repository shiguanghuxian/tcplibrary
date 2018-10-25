/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:19
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:13:07
 */
package tcplibrary

import (
	"golang.org/x/net/websocket"
)

/* websocket的连接处理，涉及到包内容解析，所以单独新建文件 */

// 用于websocket的连接处理函数
func (ws *WebSocketServer) handleConn(conn *Conn) {
	defer func() {
		if r := recover(); r != nil {
			globalLogger.Fatalf("%T", r)
		}
	}()
	// 收到消息的管道
	messageChannel := make(chan interface{}, DefaultMessageChannelSize)
	go ws.handleMessage(conn, messageChannel)
	// 循环读取 websocket
	for {
		select {
		case <-ws.ctx.Done():
			globalLogger.Infof("ws handleConn收到ctx.Done()")
			return
		default:
			// 解析websocket传输的包
			defaultPacket := new(DefaultPacket)
			err := ws.packet.GetWebsocketCodec().Receive(conn.Conn.(*websocket.Conn), defaultPacket)
			if err != nil {
				globalLogger.Errorf(err.Error())
				// 关闭连接，并通知错误
				ws.closeConn(conn, err)
				break
			}
			// 向管道写入数据
			messageChannel <- defaultPacket
		}
	}
}
