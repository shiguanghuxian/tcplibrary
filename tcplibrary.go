/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:47
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:48:57
 */
package tcplibrary

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

/* 通讯库父类 */

// 定义的tcp读缓存区大小
var (
	DefaultBufferSize         = 1024
	DefaultMessageChannelSize = 32
	DefaultReadDeadline       = 15 * time.Second
)

// TCPLibrary tcp库父类 - 每实例化
type TCPLibrary struct {
	ctx            context.Context    // 请求上下文 - 用于优雅关闭tcp服务
	cancel         context.CancelFunc // 请求上下文关闭函数
	socket         Socket             // socket 需要实现的几个方法
	packet         Packet             // 解包和封包
	readDeadline   time.Duration      // 读超时
	readBufferSize int                // 读数据时的字节缓冲
	clients        *sync.Map          // 客户端连接字典
	isServer       bool               // 是否是服务端

	listener net.Listener // tcp监听 - 只存储tcp的 不存储ws的
}

// NewTCPServer 创建TCPLibrary对象 - 只用于创建服务端对象时使用，客户端直接使用 NewTCPClient
func NewTCPServer(debug bool, socket ServerSocket, packets ...Packet) (*TCPLibrary, error) {
	if socket == nil {
		return nil, errors.New("ServerSocket参数不能是nil")
	}
	// 设置日志为debug
	if logger, ok := globalLogger.(*Logger); ok == true {
		logger.SetDefaultDebug(debug)
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
	return &TCPLibrary{
		ctx:            ctx,
		cancel:         cancel,
		packet:         packet,
		socket:         socket,
		readDeadline:   DefaultReadDeadline,
		readBufferSize: DefaultBufferSize,
		isServer:       true,
		clients:        new(sync.Map),
	}, nil
}

// GetClients 获取客户端列表,在自己的业务中使用，使用时切记小心操作
func (t *TCPLibrary) GetClients() *sync.Map {
	return t.clients
}

// SetReadDeadline 设置参数 readDeadline
func (t *TCPLibrary) SetReadDeadline(duration time.Duration) {
	t.readDeadline = duration
}

// SetReadBufferSize 设置参数 readBufferSize
func (t *TCPLibrary) SetReadBufferSize(readBufferSize int) {
	t.readBufferSize = readBufferSize
}

// StopService 停止服务
func (t *TCPLibrary) StopService() (err error) {
	// ctx 关闭 - 通知所有可以取得上下文函数结束
	t.cancel()
	// 关闭监听
	if t.listener != nil {
		err = t.listener.Close()
	}
	if err != nil {
		globalLogger.Errorf("关闭监听错误:%v", err)
		return err
	}
	// 关闭所有连接
	t.clients.Range(func(k, v interface{}) bool {
		if conn, ok := v.(*Conn); ok == true {
			t.closeConn(conn, nil)
		}
		return true
	})

	return nil
}

// 收到消息时处理
func (t *TCPLibrary) handleMessage(conn *Conn, message chan interface{}) {
	for {
		select {
		case v := <-message:
			if t.isServer == true {
				// 设置超时
				conn.SetReadDeadline(time.Now().Add(t.readDeadline))
			}
			// 调用消息回调
			go t.socket.OnRecMessage(t.ctx, conn, v)
		case <-t.ctx.Done(): // 收到结束任务消息，结束
			globalLogger.Infof("handleMessage收到ctx.Done()")
			return
		}
	}
}

// DelClients 删除一个客户端对象
func (t *TCPLibrary) delClients(keys ...interface{}) {
	for _, v := range keys {
		t.clients.Delete(v)
	}
}

// 处理从连接中读取数据
func (t *TCPLibrary) handleConn(conn *Conn) {
	defer func() {
		if r := recover(); r != nil {
			globalLogger.Fatalf("%T", r)
		}
	}()
	// 收到消息的管道
	messageChannel := make(chan interface{}, DefaultMessageChannelSize)
	go t.handleMessage(conn, messageChannel)
	// 缓冲区大小
	bufferSize := t.readBufferSize
	if bufferSize == 0 {
		bufferSize = DefaultBufferSize
	}
	data := make([]byte, 0)
	buf := make([]byte, bufferSize)
	for {
		select {
		case <-t.ctx.Done():
			globalLogger.Infof("tcp handleConn收到ctx.Done()")
			return
		default:
			n, err := conn.Read(buf)
			if err != nil {
				// 关闭连接，并通知错误
				t.closeConn(conn, err)
				return
			}
			// 解包
			data, err = conn.packet.Unmarshal(append(data, buf[:n]...), messageChannel)
			if err != nil {
				globalLogger.Errorf("%s", err.Error())
				t.socket.OnError(err)
			}
		}
	}
}

// 关闭连接，并通知错误
func (t *TCPLibrary) closeConn(conn *Conn, err error) {
	// 判断错误是不是nil和io.EOF
	if err != nil && err != io.EOF {
		globalLogger.Errorf(err.Error())
		// 通知错误
		t.socket.OnError(err)
	} else if err == io.EOF {
		err = nil
	}
	// 关闭连接
	t.socket.OnClose(conn, err) // 通知关闭
	// 删除客户端连接
	t.delClients(conn.clientID)
	err = conn.Close()
	if err != nil {
		globalLogger.Errorf("%s", err.Error())
	}
}

// CloseForClientID 根据clientID关闭连接
func (t *TCPLibrary) CloseForClientID(clientID string) error {
	// log.Println(clientID)
	connInterface, ok := t.clients.Load(clientID)
	if ok == false {
		return fmt.Errorf("踢人失败，没有这样的连接1:%s", clientID)
	}
	if conn, ok := connInterface.(*Conn); ok == true {
		t.closeConn(conn, nil)
	} else {
		return fmt.Errorf("踢人失败，没有这样的连接2:%s", clientID)
	}
	return nil
}
