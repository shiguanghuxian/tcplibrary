/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:47
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:48:57
 */
package tcplibrary

import (
	"errors"
	"io"
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

var (
	// 保存所有
	clients *sync.Map
	// 是否是服务端
	isServer bool
)

func init() {
	// 初始化客户端存储map
	clients = new(sync.Map)
	// 默认非服务端
	isServer = false
}

// GetClients 获取客户端列表,在自己的业务中使用，使用时切记小心操作
func GetClients() *sync.Map {
	return clients
}

// TCPLibrary tcp库父类
type TCPLibrary struct {
	socket         Socket        // socket 需要实现的几个方法
	packet         Packet        // 解包和封包
	readDeadline   time.Duration // 读超时
	readBufferSize int           // 读数据时的字节缓冲
}

// SetReadDeadline 设置参数 readDeadline
func (t *TCPLibrary) SetReadDeadline(duration time.Duration) {
	t.readDeadline = duration
}

// SetReadBufferSize 设置参数 readBufferSize
func (t *TCPLibrary) SetReadBufferSize(readBufferSize int) {
	t.readBufferSize = readBufferSize
}

// 收到消息时处理
func (t *TCPLibrary) handleMessage(conn *Conn, message chan interface{}) {
	for {
		select {
		case v := <-message:
			if isServer == true {
				// 设置超时
				conn.SetReadDeadline(time.Now().Add(t.readDeadline))
			}
			// 调用消息回调
			go t.socket.OnRecMessage(conn, v)
		}
	}
}

// DelClients 删除一个客户端对象
func (t *TCPLibrary) delClients(keys ...interface{}) {
	for _, v := range keys {
		clients.Delete(v)
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
		n, err := conn.Read(buf)
		if err != nil {
			// 关闭连接，并通知错误
			t.closeConn(conn, err)
			break
		}
		// 解包
		data, err = conn.packet.Unmarshal(append(data, buf[:n]...), messageChannel)
		if err != nil {
			globalLogger.Errorf("%s", err.Error())
			t.socket.OnError(err)
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
	} else {
		err = errors.New("")
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
