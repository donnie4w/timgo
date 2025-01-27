// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package timgo

import (
	"errors"
	. "github.com/donnie4w/gofer/thrift"
	. "github.com/donnie4w/timgo/stub"
)

// Register	with username and password
// domain can be set "" where is not requied；Different domains cannot communicate with each other
// 如果不需要使用 domain（域）时，可设置为空字符串，不同域无法相互通讯
func (tc *TimClient) Register(username, pwd, domain string) (ack *TimAck, err error) {
	if bs, err := sendsync(tc.syncUrl, tc.cfg.Origin, tc.ts.register(username, pwd, domain)); err == nil && bs != nil {
		ack, err = TDecode[*TimAck](bs[1:], &TimAck{})
	} else {
		err = errors.New("register failed")
	}
	return
}

// Token get a token for login
func (tc *TimClient) Token(username, pwd, domain string) (token int64, err error) {
	if bs, err := sendsync(tc.syncUrl, tc.cfg.Origin, tc.ts.token(username, pwd, domain)); err == nil && bs != nil {
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

// Login resource is the terminal information defined by the developer. For example, phone model: HUAWEI P50 Pro, iPhone 11 Pro
// if resource is not required, pass ""
// resource是开发者自定义的终端信息，一般是登录设备信息，如 HUAWEI P50 Pro，若不使用，赋空值即可
func (tc *TimClient) Login(name, pwd, domain, resource string, termtyp int8, extend map[string]string) (err error) {
	tc.closeconnect()
	tc.ts.loginByAccount(name, pwd, domain, resource, termtyp, extend)
	return tc.connect()
}

// LoginByToken login with token
func (tc *TimClient) LoginByToken(token string, resource string, termtyp int8, extend map[string]string) (err error) {
	tc.closeconnect()
	tc.ts.loginByToken(token, resource, termtyp, extend)
	return tc.connect()
}

// Logout 退出登录
func (tc *TimClient) Logout() (err error) {
	tc.ts = &tx{NewTimAuth()}
	return tc.close()
}

// ModifyPwd  password
func (tc *TimClient) ModifyPwd(oldpwd, newpwd string) (err error) {
	return tc.handler.Send(tc.ts.modifyPwd(oldpwd, newpwd))
}

// MessageToUser
// send message to a user
// showType and textType is a value defined by the developer and is sent to the peer terminal as is
// showType 和 textType 为开发者自定义字段，会原样发送到对方的终端，由开发者自定义解析，
// If showType or textType  is not required, pass 0
// showType 或 textType  不使用时，传默认值0
func (tc *TimClient) MessageToUser(user string, msg string, udshow int16, udtype int16, extend map[string]string, extra map[string][]byte) (err error) {
	return tc.handler.Send(tc.ts.message2Friend(msg, user, udshow, udtype, extend, extra))
}

// RevokeMessage
// revoke the message 撤回信息
// mid is message's id
// mid and to is  required
func (tc *TimClient) RevokeMessage(mid int64, to, room string, msg string, udshow int16, udtype int16) (err error) {
	return tc.handler.Send(tc.ts.revokeMessage(mid, to, room, msg, udshow, udtype))
}

// BurnMessage
// burn After Reading  阅后即焚
// mid is message's id
// mid and to is  required
func (tc *TimClient) BurnMessage(mid int64, to string, msg string, udshow int16, udtype int16) (err error) {
	return tc.handler.Send(tc.ts.burnMessage(mid, msg, to, udshow, udtype))
}

// MessageToRoom
// send message to a room
func (tc *TimClient) MessageToRoom(room string, msg string, udshow int16, udtype int16, extend map[string]string, extra map[string][]byte) (err error) {
	return tc.handler.Send(tc.ts.message2Room(msg, room, udshow, udtype, extend, extra))
}

// MessageByPrivacy
// send a message to a room member
func (tc *TimClient) MessageByPrivacy(user, room string, msg string, udshow int16, udtype int16, extend map[string]string, extra map[string][]byte) (err error) {
	return tc.handler.Send(tc.ts.messageByPrivacy(msg, user, room, udshow, udtype, extend, extra))
}

// StreamToUser
// send  stream data to user
func (tc *TimClient) StreamToUser(to string, msg []byte, udShow, udType int16) (err error) {
	return tc.handler.Send(tc.ts.stream(msg, to, "", udShow, udType))
}

// StreamToRoom
// send  stream data to room
func (tc *TimClient) StreamToRoom(room string, msg []byte, udShow, udType int16) (err error) {
	return tc.handler.Send(tc.ts.stream(msg, "", room, udShow, udType))
}

// PresenceToUser
// send presence to other user
// 发送状态给其他账号
func (tc *TimClient) PresenceToUser(to string, show int16, status string, subStatus int8, extend map[string]string, extra map[string][]byte) (err error) {
	return tc.handler.Send(tc.ts.presence(to, show, status, subStatus, extend, extra))
}

// PresenceToList
// send presence to other user list
func (tc *TimClient) PresenceToList(toList []string, show int16, status string, subStatus int8, extend map[string]string, extra map[string][]byte) (err error) {
	return tc.handler.Send(tc.ts.presenceList(show, status, subStatus, extend, extra, toList))
}

// BroadPresence
// broad the presence and substatus to all the friends
// 向所有好友广播状态和订阅状态
func (tc *TimClient) BroadPresence(subStatus int8, show int16, status string) (err error) {
	return tc.handler.Send(tc.ts.broadPresence(subStatus, show, status))
}

// Roster
// triggers tim to send user rosters
// 触发tim服务器发送用户花名册
func (tc *TimClient) Roster() (err error) {
	return tc.handler.Send(tc.ts.roster())
}

// Addroster  send request to the account for add friend
func (tc *TimClient) Addroster(node string, msg string) (err error) {
	return tc.handler.Send(tc.ts.addroster(node, msg))
}

// Rmroster
// remove a relationship with a specified account
// 移除与指定账号的关系
func (tc *TimClient) Rmroster(node string) (err error) {
	return tc.handler.Send(tc.ts.rmroster(node))
}

// Blockroster
// Block specified account
// 拉黑指定账号
func (tc *TimClient) Blockroster(node string) (err error) {
	return tc.handler.Send(tc.ts.blockroster(node))
}

// PullUserMessage
// pull message with user
// 拉取用户聊天消息
func (tc *TimClient) PullUserMessage(to string, mid, limit int64) (err error) {
	return tc.handler.Send(tc.ts.pullmsg(1, to, mid, limit))
}

// PullRoomMessage
// pull message of group
// 拉取群信息
func (tc *TimClient) PullRoomMessage(to string, mid, limit int64) (err error) {
	return tc.handler.Send(tc.ts.pullmsg(2, to, mid, limit))
}

// OfflineMsg
// triggers tim to send the offlien message
// 触发tim服务器发送离线信息
func (tc *TimClient) OfflineMsg() (err error) {
	return tc.handler.Send(tc.ts.offlinemsg())
}

// UserRoom
// triggers tim to send the user's ROOM account
// 触发tim服务器发送用户的群账号
func (tc *TimClient) UserRoom() (err error) {
	return tc.handler.Send(tc.ts.userroom())
}

// RoomUsers
// triggers tim to send the ROOM member account
// 触发tim服务器发送群成员账号
func (tc *TimClient) RoomUsers(node string) (err error) {
	return tc.handler.Send(tc.ts.roomusers(node))
}

// NewRoom
// creating a room, provide the room name and room type
// 创建群，需提供群名称和群类型
func (tc *TimClient) NewRoom(gtype ROOMTYPE, roomname string) (err error) {
	return tc.handler.Send(tc.ts.newroom(int64(gtype), roomname))
}

// AddRoom
// send a request to join the group
// 发送一个加入群的请求
func (tc *TimClient) AddRoom(node, msg string) (err error) {
	return tc.handler.Send(tc.ts.addroom(node, msg))
}

// PullInRoom
// pull a account into the group
// 将用户拉入群
func (tc *TimClient) PullInRoom(roomNode, userNode string) (err error) {
	return tc.handler.Send(tc.ts.pullroom(roomNode, userNode))
}

// RejectRoom
// reject a account to join into the group
// 拒绝用户加入群
func (tc *TimClient) RejectRoom(roomNode, userNode, msg string) (err error) {
	return tc.handler.Send(tc.ts.nopassroom(roomNode, userNode, msg))
}

// KickRoom
// Kick a account out of the group
// 将用户踢出群
func (tc *TimClient) KickRoom(roomNode, userNode string) (err error) {
	return tc.handler.Send(tc.ts.kickroom(roomNode, userNode))
}

// LeaveRoom
// leave group
// 退出群
func (tc *TimClient) LeaveRoom(roomNode string) (err error) {
	return tc.handler.Send(tc.ts.leaveroom(roomNode))
}

// CancelRoom
// Cancel a group
// 注销群
func (tc *TimClient) CancelRoom(roomNode string) (err error) {
	return tc.handler.Send(tc.ts.cancelroom(roomNode))
}

// BlockRoom
// block the group
// 拉黑群，拒绝被群主拉入群
func (tc *TimClient) BlockRoom(roomNode string) (err error) {
	return tc.handler.Send(tc.ts.blockroom(roomNode))
}

// BlockRoomMember
// block the group member or the account join into group
// 拉黑群成员或阻止其他账号入群
func (tc *TimClient) BlockRoomMember(roomNode, toNode string) (err error) {
	return tc.handler.Send(tc.ts.blockroomMember(roomNode, toNode))
}

// BlockRosterList
// blocklist of user
// 用户黑名单
func (tc *TimClient) BlockRosterList() (err error) {
	return tc.handler.Send(tc.ts.blockrosterlist())
}

// BlockRoomList
// blocklist of user group
// 用户群黑名单
func (tc *TimClient) BlockRoomList() (err error) {
	return tc.handler.Send(tc.ts.blockroomlist())
}

// BlockRoomMemberlist
// blocklist of group
// 群黑名单
func (tc *TimClient) BlockRoomMemberlist(node string) (err error) {
	return tc.handler.Send(tc.ts.blockroomMemberlist(node))
}

// BigDataString
// send big string
func (tc *TimClient) BigDataString(node, datastring string) (err error) {
	return tc.handler.Send(tc.ts.bigString(node, datastring))
}

// BigDataBinary
// send big binary
func (tc *TimClient) BigDataBinary(node string, dataBinary []byte) (err error) {
	return tc.handler.Send(tc.ts.bigBinary(node, dataBinary))
}

// VirtualroomRegister
// creating a Virtual room
// 创建虚拟房间
func (tc *TimClient) VirtualroomRegister() (err error) {
	return tc.handler.Send(tc.ts.virtualroom(1, "", "", 0))
}

// VirtualroomRemove
// creating a Virtual room
// 销毁虚拟房间
func (tc *TimClient) VirtualroomRemove(roomNode string) (err error) {
	return tc.handler.Send(tc.ts.virtualroom(2, roomNode, "", 0))
}

// VirtualroomAddAuth
// add push stream data permissions for virtual rooms to a account
// 给账户添加向虚拟房间推送流数据的权限
func (tc *TimClient) VirtualroomAddAuth(roomNode string, tonode string) (err error) {
	return tc.handler.Send(tc.ts.virtualroom(3, roomNode, tonode, 0))
}

// VirtualroomDelAuth
// delete the push stream data permissions for virtual rooms to a account
// 删除用户向虚拟房间推送流数据的权限
func (tc *TimClient) VirtualroomDelAuth(roomNode string, tonode string) (err error) {
	return tc.handler.Send(tc.ts.virtualroom(4, roomNode, tonode, 0))
}

// VirtualroomSub
// Subscribe the stream data of the virtual room
// 向虚拟房间订阅流数据
func (tc *TimClient) VirtualroomSub(roomNode string) (err error) {
	return tc.handler.Send(tc.ts.virtualroom(5, roomNode, "", 0))
}

// VirtualroomSubBinary
// Subscribe the stream data of the virtual room
// 向虚拟房间订阅流数据
func (tc *TimClient) VirtualroomSubBinary(roomNode string) (err error) {
	return tc.handler.Send(tc.ts.virtualroom(5, roomNode, "", 1))
}

// VirtualroomSubCancel
// cancel subscribe the stream data of the virtual room
// 取消订阅虚拟房间数据
func (tc *TimClient) VirtualroomSubCancel(roomNode string) (err error) {
	return tc.handler.Send(tc.ts.virtualroom(6, roomNode, "", 0))
}

// PushStream
// push the stream data to the virtual room
// body: body is stream data
// dtype : dtype is a data type defined by the developer and can be set to 0 if it is not required
// 推送流数据到虚拟房间
// body ：body是流数据
// dtype：dtype 是开发者自定义的数据类型，若不需要，可以设置为0
func (tc *TimClient) PushStream(virtualroom string, body []byte, dtype int8) (err error) {
	return tc.handler.Send(tc.ts.pushstream(virtualroom, body, dtype))
}

// UserInfo
// get user information
// 获取用户资料
func (tc *TimClient) UserInfo(node ...string) (err error) {
	return tc.handler.Send(tc.ts.nodeinfo(NODEINFO_USERINFO, node, nil, nil))
}

// RoomInfo
// get group information
// 获取群资料
func (tc *TimClient) RoomInfo(node ...string) (err error) {
	return tc.handler.Send(tc.ts.nodeinfo(NODEINFO_ROOMINFO, node, nil, nil))
}

// ModifyUserInfo
// modify user information
// 修改用户资料
func (tc *TimClient) ModifyUserInfo(tu *TimUserBean) (err error) {
	return tc.handler.Send(tc.ts.nodeinfo(NODEINFO_MODIFYUSER, nil, map[string]*TimUserBean{"": tu}, nil))
}

// ModifyRoomInfo
// modify group information
// 修改群资料
func (tc *TimClient) ModifyRoomInfo(roomNode string, tr *TimRoomBean) (err error) {
	return tc.handler.Send(tc.ts.nodeinfo(NODEINFO_MODIFYROOM, nil, nil, map[string]*TimRoomBean{roomNode: tr}))
}
