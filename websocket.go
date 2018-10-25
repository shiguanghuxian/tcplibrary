/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:50
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:23:42
 */
package tcplibrary

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

/* websocket 服务端 */

// WebSocketServer websocket 服务端操作对象
type WebSocketServer struct {
	*TCPLibrary
	isListener bool // 是否已监听
}

// NewWebSocketServer 创建一个websocket监听
func (t *TCPLibrary) NewWebSocketServer() (*WebSocketServer, error) {
	return &WebSocketServer{
		TCPLibrary: t,
		isListener: false,
	}, nil
}

// ListenAndServe 开始ws监听
// address 监听的地址和端口
// route 监听的路由(url)
func (ws *WebSocketServer) ListenAndServe(address, route string) error {
	if ws.isListener == true {
		return errors.New("已调用监听端口")
	}
	if address == "" {
		return errors.New("监听地址不能为空")
	}
	if route == "" {
		route = "/"
	}
	// 判断是否设置读超时
	if ws.readDeadline == 0 {
		ws.readDeadline = DefaultReadDeadline
	}
	http.Handle(route, websocket.Handler(ws.handleWebSocketConn))
	globalLogger.Infof("web socket start, net websocket addr %s", address)
	err := http.ListenAndServe(address, nil)
	return err
}

// ListenAndServeTLS 开始ws监听 tls
// address 监听的地址和端口
// route 监听的路由(url)
func (ws *WebSocketServer) ListenAndServeTLS(address, certFile, keyFile, route string) error {
	if ws.isListener == true {
		return errors.New("已调用监听端口")
	}
	if address == "" {
		return errors.New("监听地址不能为空")
	}
	if route == "" {
		route = "/"
	}
	// 判断是否设置读超时
	if ws.readDeadline == 0 {
		ws.readDeadline = DefaultReadDeadline
	}
	http.Handle(route, websocket.Handler(ws.handleWebSocketConn))
	globalLogger.Infof("web socket start, net websocket addr %s", address)
	err := http.ListenAndServeTLS(address, certFile, keyFile, nil)
	return err
}

// 处理WebSocket数据
func (ws *WebSocketServer) handleWebSocketConn(wsConn *websocket.Conn) {
	// 构建Conn对象
	conn := &Conn{
		Conn:     wsConn,
		connType: WebSocketType,
		packet:   ws.packet,
	}
	// 保存连接到客户端数组
	serverSocket, ok := ws.socket.(ServerSocket)
	if ok == false {
		// 如果建立连接函数返回false，则关闭连接
		ws.socket.OnClose(conn, fmt.Errorf("%s", "转换为ServerSocket错误")) // 通知关闭
		err := conn.Close()                                            // 关闭连接
		if err != nil {
			globalLogger.Errorf("%s", err.Error())
		}
		return
	}
	// 补上客户端id和封包解包对象，并存入服务端客户端对象
	clientID := serverSocket.GetClientID()
	conn.clientID = clientID
	ws.clients.Store(clientID, conn)
	// 设置超时
	conn.SetReadDeadline(time.Now().Add(ws.readDeadline))
	// 调用OnConnect
	// 通知连接创建后函数
	err := ws.socket.OnConnect(conn)
	if err != nil {
		// 如果建立连接函数返回false，则关闭连接
		ws.socket.OnClose(conn, err) // 通知关闭
		err = conn.Close()           // 关闭连接
		if err != nil {
			globalLogger.Errorf("%s", err.Error())
		}
		return
	}
	// 调用websocket连接处理方法
	ws.handleConn(conn)
}
