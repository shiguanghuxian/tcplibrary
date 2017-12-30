/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:15
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 14:02:48
 */
package tcplibrary

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// DefaultPacket 协议包
type DefaultPacket struct {
	Length  int32       `json:"Length"`  // Payload 包长度,4字节
	Payload interface{} `json:"Payload"` // 报文内容,n字节
}

// GetPayload 获取包内容
// 之所以有此函数，是为了兼容websocket，websocket使用此库时可直接使用json传输数据
func (dp *DefaultPacket) GetPayload() []byte {
	switch dp.Payload.(type) {
	case string:
		return []byte(dp.Payload.(string))
	case []byte:
		return dp.Payload.([]byte)
	default:
		js, err := json.Marshal(dp.Payload)
		if err == nil {
			return js
		}
		globalLogger.Errorf("默认包结构获取错误：%v", err)
	}
	return make([]byte, 0)
}

// Unmarshal 默认解包
func (dp *DefaultPacket) Unmarshal(data []byte, c chan interface{}) (outData []byte, err error) {
	// 捕获异常
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%T", r)
			globalLogger.Fatalf("默认解包错误：%v", err)
		}
	}()
	// 长度不足4个字节无法获取包长度
	if len(data) < 4 {
		return data, err
	}
	// 获取包长度
	packetLength := BytesToInt32(data[0:4]) + 4
	// 判断是否达到一个包长,没有达到直接返回
	if len(data) < int(packetLength) {
		return data, err
	}
	// 截取一个包的长度，解包
	packetData := data[:packetLength]
	// 解析内容和长度
	packet := new(DefaultPacket)
	packet.Length = int32(packetLength)
	packet.Payload = packetData[4:]
	// 写入管道数据，用于通知实际业务逻辑
	c <- packet
	// 递归调用解包
	return dp.Unmarshal(data[packetLength:], c)
}

// Marshal 默认封包
func (dp *DefaultPacket) Marshal(v interface{}) ([]byte, error) {
	packet, ok := v.(*DefaultPacket)
	if ok == false {
		return nil, errors.New("封包参数不是*DefaultPacket")
	}
	// 获取内容
	payload := packet.GetPayload()
	packet.Length = int32(len(payload))

	/* 创建Buffer对象，写入头和数据 */
	packetData := bytes.NewBuffer([]byte{})
	lengthByte := IntToBytes(packet.Length) // 长度转byte
	packetData.Write(lengthByte)
	packetData.Write(payload)

	// 返回编码后的字节数组
	return packetData.Bytes(), nil
}

// MarshalToJSON 编码到json, 同时将Payload转为字符串
func (dp *DefaultPacket) MarshalToJSON(v interface{}) ([]byte, error) {
	packet, ok := v.(*DefaultPacket)
	if ok == false {
		return nil, errors.New("封包参数不是*DefaultPacket")
	}
	packet.Payload = string(packet.GetPayload())
	// 直接转json返回
	return json.Marshal(packet)
}
