/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:41
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:22:40
 */
package tcplibrary

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"
)

/* tcp 连接 */

// TCPServer tcp服务端对象
type TCPServer struct {
	*TCPLibrary
	isListener bool // 是否已监听
}

// NewTCPServer 创建一个server实例
func (t *TCPLibrary) NewTCPServer() (*TCPServer, error) {
	return &TCPServer{
		TCPLibrary: t,
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

	return tcp.Serve(listen)
}

// ListenAndServeTLS 开始tcp监听 tls
func (tcp *TCPServer) ListenAndServeTLS(address, certFile, keyFile string) error {
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
	// 证书配置
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	tlsListener := tls.NewListener(listen, config)

	return tcp.Serve(tlsListener)
}

// Serve 开启服务
func (tcp *TCPServer) Serve(listen net.Listener) error {
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
		select {
		case <-tcp.ctx.Done():
			globalLogger.Infof("tcp Serve收到ctx.Done()")
			return nil
		default:
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
			tcp.clients.Store(clientID, conn)
			// 设置超时
			conn.SetReadDeadline(time.Now().Add(tcp.readDeadline))
			// 开启一个协程处理数据接收
			go tcp.handleConn(conn)
		}
	}
	return nil
}
