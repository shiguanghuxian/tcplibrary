# golang的tcp通讯库

使用此库你只需要创建一个server或client的结构体，实现接口，再写两句话就可以搭建起tcp通讯

## 备注
自定义协议的小伙伴如果需要websocket通讯时，请参考default_packet.go中的`GetPayload()`函数，包内容问题

当需要获取连接列表时 可以调用`GetClients()`方法获取连接`*sync.Map`
