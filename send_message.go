/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:38
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:14:28
 */
package tcplibrary

import (
	"errors"
)

/* 公共发送数据函数 */

// SendMessageToClients 发送数据给指定客户端
// 返回值，第一个值为发送成功几个连接
// 只有服务端可调用
func SendMessageToClients(v interface{}, clientIDs ...string) (sendCount int, err error) {
	if isServer == false {
		return 0, errors.New("客户端不允许调用此函数")
	}
	for _, vv := range clientIDs {
		if val, ok := clients.Load(vv); ok == true {
			if conn, ok := val.(*Conn); ok == true {
				_, err = conn.SendMessage(v)
				if err == nil {
					sendCount++
				}
			}
		}
	}
	return sendCount, err
}

// SendMessageToAll 发送给所有客户端
// 只有服务端可调用
func SendMessageToAll(v interface{}) (int, error) {
	if isServer == false {
		return 0, errors.New("客户端不允许调用此函数")
	}
	sendCount := 0
	clients.Range(func(key, val interface{}) bool {
		if conn, ok := val.(*Conn); ok == true {
			_, err := conn.SendMessage(v)
			if err != nil {
				return true
			}
			sendCount++
		}
		return true
	})
	return sendCount, nil
}
