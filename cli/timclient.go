// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package cli

import (
	"errors"

	. "github.com/donnie4w/gofer/thrift"
	. "github.com/donnie4w/timgo/stub"
)

// register with username and password
// domain can be set "" where is not requied；Different domains cannot communicate with each other
// 如果不需要使用 domain（域）时，可设置为空字符串，不同域无法相互通讯
func (this *TimClient) Register(username, pwd, domain string) (ack *TimAck, err error) {
	if bs, err := sendsync(this.conf, this.Tx.register(username, pwd, domain)); err == nil && bs != nil {
		ack, err = TDecode[*TimAck](bs[1:], &TimAck{})
	} else {
		err = errors.New("register failed")
	}
	return
}

// get a token for login
func (this *TimClient) Token(username, pwd, domain string) (token int64, err error) {
	if bs, err := sendsync(this.conf, this.Tx.token(username, pwd, domain)); err == nil && bs != nil {
		if ta, err := TDecode[*TimAck](bs[1:], &TimAck{}); err == nil {
			if ok := ta.Ok; !ok {
				err = errors.New(*ta.Error.Info)
			} else {
				token = *ta.T
			}
		}
	} else {
		err = errors.New("get token failed")
	}
	return
}

// resource is the terminal information defined by the developer. For example, phone model: HUAWEI P50 Pro, iPhone 11 Pro
// if resource is not required, pass ""
// resource是开发者自定义的终端信息，一般是登录设备信息，如 HUAWEI P50 Pro，若不使用，赋空值即可
func (this *TimClient) Login(name, pwd, domain, resource string, termtyp int8, extend map[string]string) (err error) {
	this.close()
	this.Tx.loginByAccount(name, pwd, domain, resource, termtyp, extend)
	return this.connect()
}

func (this *TimClient) LoginByToken(token int64, resource string, termtyp int8, extend map[string]string) (err error) {
	this.close()
	this.Tx.loginByToken(token, resource, termtyp, extend)
	return this.connect()
}

// 退出登录
func (this *TimClient) Logout() (err error) {
	return this.CloseAndLogout()
}

// modify password
func (this *TimClient) Modify(oldpwd, newpwd string) (err error) {
	return this.handler.sendws(this.Tx.modify(oldpwd, newpwd))
}

// send message to a user
// showType and textType is a value defined by the developer and is sent to the peer terminal as is
// showType 和 textType 为开发者自定义字段，会原样发送到对方的终端，由开发者自定义解析，
// If showType or textType  is not required, pass 0
// showType 或 textType  不使用时，传默认值0
func (this *TimClient) MessageToUser(user string, msg string, udshow int16, udtype int16, extend map[string]string, extra map[string][]byte) (err error) {
	return this.handler.sendws(this.Tx.message2Friend(msg, user, udshow, udtype, extend, extra))
}

// revoke the message 撤回信息
// mid is message's id
// mid and to is  required
func (this *TimClient) RevokeMessage(mid int64, to, room string, msg string, udshow int16, udtype int16) (err error) {
	return this.handler.sendws(this.Tx.revokeMessage(mid, to, room, msg, udshow, udtype))
}

// Burn After Reading  阅后即焚
// mid is message's id
// mid and to is  required
func (this *TimClient) BurnMessage(mid int64, to string, msg string, udshow int16, udtype int16) (err error) {
	return this.handler.sendws(this.Tx.burnMessage(mid, msg, to, udshow, udtype))
}

// send message to a room
func (this *TimClient) MessageToRoom(room string, msg string, udshow int16, udtype int16, extend map[string]string, extra map[string][]byte) (err error) {
	return this.handler.sendws(this.Tx.message2Room(msg, room, udshow, udtype, extend, extra))
}

// send a message to a room member
func (this *TimClient) MessageByPrivacy(user, room string, msg string, udshow int16, udtype int16, extend map[string]string, extra map[string][]byte) (err error) {
	return this.handler.sendws(this.Tx.messageByPrivacy(msg, user, room, udshow, udtype, extend, extra))
}

// send  stream data to user
func (this *TimClient) StreamToUser(to string, msg []byte, udShow, udType int16) (err error) {
	return this.handler.sendws(this.Tx.stream(msg, to, "", udShow, udType))
}

// send  stream data to room
func (this *TimClient) StreamToRoom(room string, msg []byte, udShow, udType int16) (err error) {
	return this.handler.sendws(this.Tx.stream(msg, "", room, udShow, udType))
}

// send presence to other user
// 发送状态给其他账号
func (this *TimClient) PresenceToUser(to string, show int16, status string, subStatus int8, extend map[string]string, extra map[string][]byte) (err error) {
	return this.handler.sendws(this.Tx.presence(to, show, status, subStatus, extend, extra))
}

// send presence to other user list
func (this *TimClient) PresenceToList(toList []string, show int16, status string, subStatus int8, extend map[string]string, extra map[string][]byte) (err error) {
	return this.handler.sendws(this.Tx.presenceList(show, status, subStatus, extend, extra, toList))
}

// broad the presence and substatus to all the friends
// 向所有好友广播状态和订阅状态
func (this *TimClient) BroadPresence(subStatus int8, show int16, status string) (err error) {
	return this.handler.sendws(this.Tx.broadPresence(subStatus, show, status))
}

// triggers tim to send user rosters
// 触发tim服务器发送用户花名册
func (this *TimClient) Roster() (err error) {
	return this.handler.sendws(this.Tx.roster())
}

// send request to  the account for add friend
func (this *TimClient) Addroster(node string, msg string) (err error) {
	return this.handler.sendws(this.Tx.addroster(node, msg))
}

// remove a relationship with a specified account
// 移除与指定账号的关系
func (this *TimClient) Rmroster(node string) (err error) {
	return this.handler.sendws(this.Tx.rmroster(node))
}

// Block specified account
// 拉黑指定账号
func (this *TimClient) Blockroster(node string) (err error) {
	return this.handler.sendws(this.Tx.blockroster(node))
}

// pull message with user
// 拉取用户聊天消息
func (this *TimClient) PullUserMessage(to string, mid, limit int64) (err error) {
	return this.handler.sendws(this.Tx.pullmsg(1, to, mid, limit))
}

// pull message of group
// 拉取群信息
func (this *TimClient) PullRoomMessage(to string, mid, limit int64) (err error) {
	return this.handler.sendws(this.Tx.pullmsg(2, to, mid, limit))
}

// triggers tim to send the offlien message
// 触发tim服务器发送离线信息
func (this *TimClient) OfflineMsg() (err error) {
	return this.handler.sendws(this.Tx.offlinemsg())
}

// triggers tim to send the user's ROOM account
// 触发tim服务器发送用户的群账号
func (this *TimClient) UserRoom() (err error) {
	return this.handler.sendws(this.Tx.userroom())
}

// triggers tim to send the ROOM member account
// 触发tim服务器发送群成员账号
func (this *TimClient) RoomUsers(node string) (err error) {
	return this.handler.sendws(this.Tx.roomusers(node))
}

// creating a room, provide the room name and room type
// 创建群，需提供群名称和群类型
func (this *TimClient) NewRoom(gtype ROOMTYPE, roomname string) (err error) {
	return this.handler.sendws(this.Tx.newroom(int64(gtype), roomname))
}

// send a request to join the group
// 发送一个加入群的请求
func (this *TimClient) AddRoom(node, msg string) (err error) {
	return this.handler.sendws(this.Tx.addroom(node, msg))
}

// pull a account into the group
// 将用户拉入群
func (this *TimClient) PullInRoom(roomNode, userNode string) (err error) {
	return this.handler.sendws(this.Tx.pullroom(roomNode, userNode))
}

// reject a account to join into the group
// 拒绝用户加入群
func (this *TimClient) RejectRoom(roomNode, userNode, msg string) (err error) {
	return this.handler.sendws(this.Tx.nopassroom(roomNode, userNode, msg))
}

// Kick a account out of the group
// 将用户踢出群
func (this *TimClient) KickRoom(roomNode, userNode string) (err error) {
	return this.handler.sendws(this.Tx.kickroom(roomNode, userNode))
}

// leave group
// 退出群
func (this *TimClient) LeaveRoom(roomNode string) (err error) {
	return this.handler.sendws(this.Tx.leaveroom(roomNode))
}

// Cancel a group
// 注销群
func (this *TimClient) CancelRoom(roomNode string) (err error) {
	return this.handler.sendws(this.Tx.cancelroom(roomNode))
}

// block the group
// 拉黑群
func (this *TimClient) BlockRoom(roomNode string) (err error) {
	return this.handler.sendws(this.Tx.blockroom(roomNode))
}

// block the group member or the account join into group
// 拉黑群成员或其他账号入群
func (this *TimClient) BlockRoomMember(roomNode, toNode string) (err error) {
	return this.handler.sendws(this.Tx.blockroomMember(roomNode, toNode))
}

// blocklist of user
// 用户黑名单
func (this *TimClient) BlockRosterList() (err error) {
	return this.handler.sendws(this.Tx.blockrosterlist())
}

// blocklist of user group
// 用户群黑名单
func (this *TimClient) BlockRoomList() (err error) {
	return this.handler.sendws(this.Tx.blockroomlist())
}

// blocklist of group
// 群黑名单
func (this *TimClient) BlockRoomMemberlist(node string) (err error) {
	return this.handler.sendws(this.Tx.blockroomMemberlist(node))
}

// creating a Virtual room
// 创建虚拟房间
func (this *TimClient) VirtualroomRegister() (err error) {
	return this.handler.sendws(this.Tx.virtualroom(1, "", ""))
}

// creating a Virtual room
// 销毁虚拟房间
func (this *TimClient) VirtualroomRemove(roomNode string) (err error) {
	return this.handler.sendws(this.Tx.virtualroom(2, roomNode, ""))
}

// Add push stream data permissions for virtual rooms to a account
// 给账户添加向虚拟房间推送流数据的权限
func (this *TimClient) VirtualroomAddAuth(roomNode string, tonode string) (err error) {
	return this.handler.sendws(this.Tx.virtualroom(3, roomNode, tonode))
}

// delete the push stream data permissions for virtual rooms to a account
// 删除用户向虚拟房间推送流数据的权限
func (this *TimClient) VirtualroomDelAuth(roomNode string, tonode string) (err error) {
	return this.handler.sendws(this.Tx.virtualroom(4, roomNode, tonode))
}

// Subscribe the stream data of the virtual room
// 向虚拟房间订阅流数据
func (this *TimClient) VirtualroomSub(roomNode string) (err error) {
	return this.handler.sendws(this.Tx.virtualroom(5, roomNode, ""))
}

// cancel subscribe the stream data of the virtual room
// 取消订阅虚拟房间数据
func (this *TimClient) VirtualroomSubCancel(roomNode string) (err error) {
	return this.handler.sendws(this.Tx.virtualroom(6, roomNode, ""))
}

// push the stream data to the virtual room
// body: body is stream data
// dtype : dtype is a data type defined by the developer and can be set to 0 if it is not required
// 推送流数据到虚拟房间
// body ：body是流数据
// dtype：dtype 是开发者自定义的数据类型，若不需要，可以设置为0
func (this *TimClient) PushStream(virtualroom string, body []byte, dtype int8) (err error) {
	return this.handler.sendws(this.Tx.pushstream(virtualroom, body, dtype))
}

// get user information
// 获取用户资料
func (this *TimClient) UserInfo(node ...string) (err error) {
	return this.handler.sendws(this.Tx.nodeinfo(NODEINFO_USERINFO, node, nil, nil))
}

// get group information
// 获取群资料
func (this *TimClient) RoomInfo(node ...string) (err error) {
	return this.handler.sendws(this.Tx.nodeinfo(NODEINFO_ROOMINFO, node, nil, nil))
}

// modify user information
// 修改用户资料
func (this *TimClient) ModifyUserInfo(tu *TimUserBean) (err error) {
	return this.handler.sendws(this.Tx.nodeinfo(NODEINFO_MODIFYUSER, nil, map[string]*TimUserBean{"": tu}, nil))
}

// modify group information
// 修改群资料
func (this *TimClient) ModifyRoomInfo(roomNode string, tr *TimRoomBean) (err error) {
	return this.handler.sendws(this.Tx.nodeinfo(NODEINFO_MODIFYROOM, nil, nil, map[string]*TimRoomBean{roomNode: tr}))
}
