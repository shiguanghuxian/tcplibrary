/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:41
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:22:40
 */
package tcplibrary

import (
	"errors"
	"fmt"
	"net"
	"time"
)

/* tcp 连接 */

// TCPServer tcp服务端对象
type TCPServer struct {
	*TCPLibrary
	listener   *net.TCPListener // tcp监听
	isListener bool             // 是否已监听
}

// NewTCPServer 创建一个server实例
func NewTCPServer(debug bool, socket ServerSocket, packets ...Packet) (*TCPServer, error) {
	if socket == nil {
		return nil, errors.New("ServerSocket参数不能是nil")
	}
	// 封包解包对象
	var packet Packet
	if len(packets) == 0 {
		packet = new(DefaultPacket)
	} else {
		packet = packets[0]
	}
	// 标记为服务端
	isServer = true

	return &TCPServer{
		TCPLibrary: &TCPLibrary{
			packet:         packet,
			socket:         socket,
			readDeadline:   DefaultReadDeadline,
			readBufferSize: DefaultBufferSize,
		},
		isListener: false,
	}, nil
}

// ListenAndServe 开始tcp监听
func (tcp *TCPServer) ListenAndServe(address string) error {
	if tcp.isListener == true {
		return errors.New("已调用监听端口")
	}
	if address == "" {
		return errors.New("监听地址不能为空")
	}
	// 开启tcp监听
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		globalLogger.Errorf(err.Error())
		return err
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	// 判断是否设置读超时
	if tcp.readDeadline == 0 {
		tcp.readDeadline = DefaultReadDeadline
	}
	// 监听对象赋值给当前对象，并将isListener设为true
	tcp.listener = listen
	tcp.isListener = true
	// 打印开启tcp服务
	globalLogger.Infof("tcp socket start, net %s addr %s", listen.Addr().Network(), listen.Addr().String())
	// 开始接收客户端连接
	for {
		tcpConn, err := tcp.listener.Accept()
		if err != nil {
			globalLogger.Errorf(err.Error())
			continue
		}
		// 创建一个Conn对象
		conn := &Conn{
			Conn:     tcpConn,
			connType: TCPSocketType,
			packet:   tcp.packet,
		}
		// 获取客户端id
		serverSocket, ok := tcp.socket.(ServerSocket)
		if ok == false {
			// 如果建立连接函数返回false，则关闭连接
			tcp.socket.OnClose(conn, fmt.Errorf("%s", "转换为ServerSocket错误")) // 通知关闭
			err = conn.Close()                                              // 关闭连接
			if err != nil {
				globalLogger.Errorf(err.Error())
			}
			break
		}
		clientID := serverSocket.GetClientID()
		conn.clientID = clientID
		// 通知连接创建后函数
		err = tcp.socket.OnConnect(conn)
		if err != nil {
			// 如果建立连接函数返回false，则关闭连接
			tcp.socket.OnClose(conn, err) // 通知关闭
			err = conn.Close()            // 关闭连接
			if err != nil {
				globalLogger.Errorf(err.Error())
			}
			break
		}
		clients.Store(clientID, conn)
		// 设置超时
		conn.SetReadDeadline(time.Now().Add(tcp.readDeadline))
		// 开启一个协程处理数据接收
		go tcp.handleConn(conn)
	}
	return nil
}
