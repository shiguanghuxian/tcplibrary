package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/shiguanghuxian/tcplibrary"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	for i := 0; i < 100; i++ {
		go func() {
			client := new(Client)
			c, err := tcplibrary.NewTCPClient(true, client)
			if err != nil {
				log.Println(err)
			}
			err = c.DialAndStart(":1028")
			log.Println(err)
		}()
	}

	select {}
}

// Client tcp客户端服务对象
type Client struct {
}

// OnConnect 连接建立时
func (c *Client) OnConnect(conn *tcplibrary.Conn) error {
	log.Println("OnConnect")
	go func() {
		conn := conn
		for {
			pp := &tcplibrary.DefaultPacket{
				Payload: []byte("你好世界"),
			}
			n, err := conn.SendMessage(pp)
			log.Println(n, err)
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}

// OnError 遇到错误时
func (c *Client) OnError(err error) {
	log.Println("OnError")
	log.Println(err)
}

// OnClose 连接关闭时
func (c *Client) OnClose(conn *tcplibrary.Conn, err error) {
	log.Println("OnClose")
	log.Println(err)
	os.Exit(1)
}

// OnRecMessage 收到消息时
func (c *Client) OnRecMessage(conn *tcplibrary.Conn, v interface{}) {
	log.Println("OnRecMessage")
	if packet, ok := v.(*tcplibrary.DefaultPacket); ok == true {
		log.Printf("消息体长度:%d 消息体内容:%s", packet.Length, string(packet.GetPayload()))
	} else {
		js, _ := json.Marshal(v)
		log.Println(string(js))
	}
}
