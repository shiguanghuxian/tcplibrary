package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/shiguanghuxian/tcplibrary"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	server := new(Server)
	// 启动websocket监听
	webSocketServer, err := tcplibrary.NewWebSocketServer(true, server)
	if err != nil {
		log.Println(err)
	}
	go func() {
		err = webSocketServer.ListenAndServe(":1126", "/vivi")
		if err != nil {
			log.Println(err)
		}
	}()

	// 启动tcp监听
	tcpServer, err := tcplibrary.NewTCPServer(true, server)
	if err != nil {
		log.Println(err)
	}
	err = tcpServer.ListenAndServe(":1028")
	// 启动
	log.Println(err)
}

// Server 服务端服务对象
type Server struct {
}

// OnConnect 连接建立时
func (s *Server) OnConnect(conn *tcplibrary.Conn) error {
	log.Println("OnConnect")
	return nil
}

// OnError 连接遇到错误时
func (s *Server) OnError(err error) {
	log.Println("OnError")
	log.Println(err)
}

// OnClose 连接关闭时
func (s *Server) OnClose(conn *tcplibrary.Conn, err error) {
	log.Println("OnClose")
	log.Println(err)
}

// OnRecMessage 收到客户端发送过来的消息时
func (s *Server) OnRecMessage(conn *tcplibrary.Conn, v interface{}) {
	log.Println("OnRecMessage")
	if packet, ok := v.(*tcplibrary.DefaultPacket); ok == true {
		log.Printf("消息体长度:%d 消息体内容:%s\n", packet.Length, string(packet.GetPayload()))
		// 转发给所有
		n, err := tcplibrary.SendMessageToAll(v)
		log.Printf("成功发送%d个客户端，错误:%v\n", n, err)
	} else {
		js, _ := json.Marshal(v)
		log.Println(string(js))
	}

}

var i = 0

// GetClientID 生成一个客户端连接，只要唯一即可
func (s *Server) GetClientID() string {
	i++
	return fmt.Sprint(i)
}
