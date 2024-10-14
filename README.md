## tim的go客户端

######  使用 timgo 连接 **[tim](https://github.com/donnie4w/tim)** 服务器，可以进行登录，注册，加好友，建群，发信，广播状态，发送视频流数据等操作，实现im的具体功能
######  功能具体可以参考 [webtim在线测试项目](https://tim.tlnet.top)

### 快速使用

```bash
go get github.com/donnie4w/timgo
```

```go
//建立操作对象
tc := timgo.NewTimClient("192.168.2.11", "5080", false)

//实现消息接收后的处理事件
tc.MessageHandler(func(tm *TimMessage) {})

//实现状态信息接收后的处理事件
tc.PresenceHandler(func(tp *TimPresence) {})
```

##### 注册

```go
if ack, _ := tc.Register("13912345678", "123456", "yourdomain"); ack != nil {
    fmt.Println(ack)
}
```

##### 登录

```go
tc.Login("13912345678", "123456", "yourdomain", "android 13 huawei", 1, nil)
```

##### 向好友发信

```go
tc.MessageToUser("QdH6CCms5FV", "hello", 0, 0, nil, nil)
```