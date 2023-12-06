// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package cli

import (
	. "github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/gofer/thrift"
	. "github.com/donnie4w/timgo/stub"
)

// tim 接口类型 type of interface
type TIMTYPE byte

// 群类型
type ROOMTYPE int64

const (
	TIMACK           TIMTYPE = 12 //回复ACK
	TIMPING          TIMTYPE = 13 //ping
	TIMREGISTER      TIMTYPE = 14 //注册
	TIMTOKEN         TIMTYPE = 15 //拉取token
	TIMAUTH          TIMTYPE = 16 //登录验证
	TIMOFFLINEMSG    TIMTYPE = 17 //推送离线消息
	TIMOFFLINEMSGEND TIMTYPE = 18 //推送离线消息完毕标识
	TIMBROADPRESENCE TIMTYPE = 19 //广播状态订阅信息给在线好友
	TIMLOGOUT        TIMTYPE = 20 //强制下线
	TIMPULLMESSAGE   TIMTYPE = 21 //拉取聊天信息
	TIMVROOM         TIMTYPE = 22 //虚拟房间操作
	TIMBUSINESS      TIMTYPE = 41 //业务
	TIMNODES         TIMTYPE = 42 //账号信息
	TIMMESSAGE       TIMTYPE = 90 //消息
	TIMPRESENCE      TIMTYPE = 91 //状态
	TIMREVOKEMESSAGE TIMTYPE = 92 //撤回
	TIMBURNMESSAGE   TIMTYPE = 93 //焚烧信息
	TIMSTREAM        TIMTYPE = 94 //虚拟房间流数据
	//room species 群种类
	ROOM_PRIVATE ROOMTYPE = 1 //私有群，入群需验证
	ROOM_OPEN    ROOMTYPE = 2 //公开群，入群不需验证
)

var (
	BUSINESS_ROSTER              int32 = 1  //拉取花名册
	BUSINESS_USERROOM            int32 = 2  //拉取群账号
	BUSINESS_ROOMUSERS           int32 = 3  //拉取群成员账号
	BUSINESS_ADDROSTER           int32 = 4  //加好友
	BUSINESS_FRIEND              int32 = 5  //成为好友
	BUSINESS_REMOVEROSTER        int32 = 6  //删除好友
	BUSINESS_BLOCKROSTER         int32 = 7  //拉黑账号
	BUSINESS_NEWROOM             int32 = 8  //建群
	BUSINESS_ADDROOM             int32 = 9  //加入群
	BUSINESS_PASSROOM            int32 = 10 //通过加群申请
	BUSINESS_NOPASSROOM          int32 = 11 //不通过加群申请
	BUSINESS_PULLROOM            int32 = 12 //拉人入群
	BUSINESS_KICKROOM            int32 = 13 //踢人出群
	BUSINESS_BLOCKROOM           int32 = 14 //拉黑群
	BUSINESS_BLOCKROOMMEMBER     int32 = 15 //拉黑群成员
	BUSINESS_LEAVEROOM           int32 = 16 //退群
	BUSINESS_CANCELROOM          int32 = 17 //注销群
	BUSINESS_BLOCKROSTERLIST     int32 = 18 //拉取账号黑名单
	BUSINESS_BLOCKROOMLIST       int32 = 19 //拉取账号拉黑群名单
	BUSINESS_BLOCKROOMMEMBERLIST int32 = 20 //管理员拉取群黑名单
	BUSINESS_MODIFYAUTH          int32 = 21 //修改用户密码
)

const (
	NODEINFO_ROSTER              int32 = 1  //花名册
	NODEINFO_ROOM                int32 = 2  //用户的群
	NODEINFO_ROOMMEMBER          int32 = 3  //群的成员
	NODEINFO_USERINFO            int32 = 4  //用户信息
	NODEINFO_ROOMINFO            int32 = 5  //群信息
	NODEINFO_MODIFYUSER          int32 = 6  //修改用户信息
	NODEINFO_MODIFYROOM          int32 = 7  //修改群信息
	NODEINFO_BLOCKROSTERLIST     int32 = 8  //用户黑名单
	NODEINFO_BLOCKROOMLIST       int32 = 9  //用户拉黑群的群账号
	NODEINFO_BLOCKROOMMEMBERLIST int32 = 10 //群拉黑账号名单
)

const (
	VRITURLROOM_REGISTER  int64 = 1 //虚拟房间注册
	VRITURLROOM_REMOVE    int64 = 2 //虚拟房间删除
	VRITURLROOM_ADDAUTH   int64 = 3 //虚拟房间加权
	VRITURLROOM_RMAUTH    int64 = 4 //虚拟房间除权
	VRITURLROOM_SUB       int64 = 5 //虚拟房间订阅
	VRITURLROOM_SUBCANCEL int64 = 6 //虚拟房间取消订阅
)

const (
	ERR_HASEXIST   = 4101
	ERR_NOPASS     = 4102
	ERR_EXPIREOP   = 4103
	ERR_PARAMS     = 4104
	ERR_AUTH       = 4105
	ERR_ACCOUNT    = 4106
	ERR_INTERFACE  = 4107
	ERR_CANCEL     = 4108
	ERR_NOEXIST    = 4109
	ERR_BLOCK      = 4110
	ERR_OVERENTRY  = 4111
	ERR_MODIFYAUTH = 4112
)

// 自定义订阅状态，根据实际需求赋值
const (
	SUBSTATUS_REQ int8 = 1
	SUBSTATUS_ACK int8 = 2
)

type tx struct {
	ta *TimAuth
}

func (this *tx) ping() []byte {
	buf := NewBuffer()
	buf.WriteByte(byte(TIMPING))
	return buf.Bytes()
}

// register with username password and domain
// 使用用户名密码域名注册
func (this *tx) register(username, pwd, domain string) []byte {
	ta := &TimAuth{Name: &username, Pwd: &pwd, Domain: &domain}
	buf := NewBuffer()
	buf.WriteByte(byte(TIMREGISTER))
	buf.Write(thrift.TEncode(ta))
	return buf.Bytes()
}

// get token by username password and domain
// 使用用户名密码域名获取token
func (this *tx) token(username, pwd, domain string) []byte {
	ta := &TimAuth{Name: &username, Pwd: &pwd, Domain: &domain}
	buf := NewBuffer()
	buf.WriteByte(byte(TIMTOKEN))
	buf.Write(thrift.TEncode(ta))
	return buf.Bytes()
}

// login with username password and domain
// 使用用户名密码域名登录
func (this *tx) loginByAccount(username, pwd, domain, resource string, termtyp int8, extend map[string]string) []byte {
	this.ta.Name, this.ta.Pwd, this.ta.Domain, this.ta.Resource, this.ta.Termtyp, this.ta.Extend = &username, &pwd, &domain, &resource, &termtyp, extend
	buf := NewBuffer()
	buf.WriteByte(byte(TIMAUTH))
	buf.Write(thrift.TEncode(this.ta))
	return buf.Bytes()
}

// login with token
// 使用token登录
func (this *tx) loginByToken(token int64, resource string, termtyp int8, extend map[string]string) []byte {
	this.ta.Token, this.ta.Resource, this.ta.Termtyp, this.ta.Extend = &token, &resource, &termtyp, extend
	buf := NewBuffer()
	buf.WriteByte(byte(TIMAUTH))
	buf.Write(thrift.TEncode(this.ta))
	return buf.Bytes()
}

func (this *tx) login() []byte {
	buf := NewBuffer()
	buf.WriteByte(byte(TIMAUTH))
	buf.Write(thrift.TEncode(this.ta))
	return buf.Bytes()
}

func (this *tx) _message(timtype TIMTYPE, mstype, odType int8, msg string, to string, roomId string, udshow int16, udtype int16, msgId int64, dataBinary []byte, extend map[string]string, extra map[string][]byte) []byte {
	tm := &TimMessage{MsType: mstype, OdType: odType, Extend: extend, Extra: extra, DataBinary: dataBinary}
	if roomId != "" {
		tm.RoomTid = &Tid{Node: roomId}
		if to == "" {
			tm.MsType = 3
		}
	}
	if to != "" {
		tm.ToTid = &Tid{Node: to}
	}
	if udshow > 0 {
		tm.Udshow = &udshow
	}
	if udtype > 0 {
		tm.Udtype = &udtype
	}
	if msg != "" {
		tm.DataString = &msg
	}
	if msgId > 0 {
		tm.Mid = &msgId
	}
	buf := NewBuffer()
	buf.WriteByte(byte(timtype))
	buf.Write(thrift.TEncode(tm))
	return buf.Bytes()
}

// Text friends regularly 常规发信息给朋友
// msg:Sending content 发送的信息
// to: Friend account 好友账号
// private:Whether to send private messages to room friends 是否私信给群友
// roomId:ROOM account	群账号
// udshow:Developer custom, will be handled in the other party 开发者自定义,将会在对方处理
// udtype:Like udshow 如同udshow
// extend & extra :Extended fields, developer custom  扩展字段，开发者自定义, 赋值后送到对方处理
func (this *tx) message2Friend(msg string, to string, udshow int16, udtype int16, extend map[string]string, extra map[string][]byte) []byte {
	return this._message(TIMMESSAGE, 2, 1, msg, to, "", udshow, udtype, 0, nil, extend, extra)
}

func (this *tx) messageByPrivacy(msg string, to string, roomId string, udshow int16, udtype int16, extend map[string]string, extra map[string][]byte) []byte {
	return this._message(TIMMESSAGE, 2, 1, msg, to, roomId, udshow, udtype, 0, nil, extend, extra)
}

func (this *tx) message2Room(msg string, roomId string, udshow int16, udtype int16, extend map[string]string, extra map[string][]byte) []byte {
	return this._message(TIMMESSAGE, 3, 1, msg, "", roomId, udshow, udtype, 0, nil, extend, extra)
}

func (this *tx) revokeMessage(msgId int64, to, room string, msg string, udshow int16, udtype int16) []byte {
	return this._message(TIMREVOKEMESSAGE, 2, 2, msg, to, room, udshow, udtype, msgId, nil, nil, nil)
}

func (this *tx) burnMessage(msgId int64, msg string, to string, udshow int16, udtype int16) []byte {
	return this._message(TIMBURNMESSAGE, 2, 3, msg, to, "", udshow, udtype, msgId, nil, nil, nil)
}

func (this *tx) stream(msg []byte, to string, room string, udShow, udType int16) []byte {
	tm := &TimMessage{MsType: 2, OdType: 5}
	if room != "" {
		tm.RoomTid = &Tid{Node: room}
		if to == "" {
			tm.MsType = 3
		}
	}
	if to != "" {
		tm.ToTid = &Tid{Node: to}
	}
	if udShow > 0 {
		tm.Udshow = &udShow
	}
	if udType > 0 {
		tm.Udtype = &udType
	}
	tm.DataBinary = msg
	buf := NewBuffer()
	buf.WriteByte(byte(TIMMESSAGE))
	buf.Write(thrift.TEncode(tm))
	return buf.Bytes()
}

// Extended send message 扩展发信,
// The information will reach the target user, but it will not be saved 信息会到达目标用户，但不会被保存起来
// func (this *tx) exMessage(odType int8, msg string, to string, roomId string, udshow int16, udtype int16, extend map[string]string, extra map[string][]byte) (_r []byte) {
// 	if odType > 30 {
// 		_r = this._message(TIMMESSAGE, 2, odType, msg, to, roomId, udshow, udtype, 0, extend, extra)
// 	} else {
// 		logging.Error(odType, " cannot pass,The value of odType must be greater than 30")
// 	}
// 	return
// }

func (this *tx) _presence(timtype TIMTYPE, to string, show int16, status string, toList []string, subStatus int8, extend map[string]string, extra map[string][]byte) []byte {
	tp := &TimPresence{Extend: extend, Extra: extra, ToList: toList}
	if to != "" {
		tp.ToTid = &Tid{Node: to}
	}
	if show > 0 {
		tp.Show = &show
	}
	if status != "" {
		tp.Status = &status
	}
	if subStatus > 0 {
		tp.SubStatus = &subStatus
	}
	buf := NewBuffer()
	buf.WriteByte(byte(timtype))
	buf.Write(thrift.TEncode(tp))
	return buf.Bytes()
}

func (this *tx) presence(to string, show int16, status string, subStatus int8, extend map[string]string, extra map[string][]byte) []byte {
	return this._presence(TIMPRESENCE, to, show, status, nil, subStatus, extend, extra)
}

func (this *tx) presenceList(show int16, status string, subStatus int8, extend map[string]string, extra map[string][]byte, toList []string) []byte {
	return this._presence(TIMPRESENCE, "", show, status, toList, subStatus, extend, extra)
}

func (this *tx) broadPresence(subStatus int8, show int16, status string) []byte {
	return this._presence(TIMBROADPRESENCE, "", show, status, nil, subStatus, nil, nil)
}

func (this *tx) pullmsg(rtype int32, to string, mid, limit int64) []byte {
	rt := &TimReq{Rtype: &rtype, Node: &to, ReqInt: &mid, ReqInt2: &limit}
	buf := NewBuffer()
	buf.WriteByte(byte(TIMPULLMESSAGE))
	buf.Write(thrift.TEncode(rt))
	return buf.Bytes()
}

func (this *tx) offlinemsg() []byte {
	buf := NewBuffer()
	buf.WriteByte(byte(TIMOFFLINEMSG))
	return buf.Bytes()
}

func (this *tx) addroster(node string, msg string) (_r []byte) {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_ADDROSTER, Node: &node, ReqStr: &msg}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) rmroster(node string) (_r []byte) {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_REMOVEROSTER, Node: &node}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) blockroster(node string) (_r []byte) {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_BLOCKROSTER, Node: &node}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) roster() (_r []byte) {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_ROSTER}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) modify(oldpwd, newpwd string) (_r []byte) {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_MODIFYAUTH, ReqStr: &oldpwd, ReqStr2: &newpwd}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) userroom() []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_USERROOM}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) roomusers(node string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_ROOMUSERS, Node: &node}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) newroom(gtype int64, roomname string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_NEWROOM, Node: &roomname, ReqInt: &gtype}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) addroom(gnode string, msg string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_ADDROOM, Node: &gnode, ReqStr: &msg}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) pullroom(rnode, unode string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_PULLROOM, Node: &rnode, Node2: &unode}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) nopassroom(rnode, unode, msg string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_NOPASSROOM, Node: &rnode, Node2: &unode, ReqStr: &msg}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) kickroom(rnode, unode string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_KICKROOM, Node: &rnode, Node2: &unode}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) leaveroom(gnode string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_LEAVEROOM, Node: &gnode}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) cancelroom(gnode string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_CANCELROOM, Node: &gnode}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}
func (this *tx) blockroom(gnode string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_BLOCKROOM, Node: &gnode}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) blockroomMember(gnode, tonode string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_BLOCKROOMMEMBER, Node: &gnode, Node2: &tonode}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) blockrosterlist() []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_BLOCKROSTERLIST}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) blockroomlist() []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_BLOCKROOMLIST}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) blockroomMemberlist(gnode string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &BUSINESS_BLOCKROOMMEMBERLIST, Node: &gnode}
	buf.WriteByte(byte(TIMBUSINESS))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) virtualroom(rtype int32, vnode, unode string) []byte {
	buf := NewBuffer()
	tr := &TimReq{Rtype: &rtype}
	if vnode != "" {
		tr.Node = &vnode
	}
	if unode != "" {
		tr.Node2 = &unode
	}
	buf.WriteByte(byte(TIMVROOM))
	buf.Write(thrift.TEncode(tr))
	return buf.Bytes()
}

func (this *tx) pushstream(node string, body []byte, dtype int8) []byte {
	ts := &TimStream{VNode: node, Body: body}
	if dtype != 0 {
		ts.Dtype = &dtype
	}
	buf := NewBuffer()
	buf.WriteByte(byte(TIMSTREAM))
	buf.Write(thrift.TEncode(ts))
	return buf.Bytes()
}

func (this *tx) nodeinfo(ntype int32, nodelist []string, usermap map[string]*TimUserBean, roommap map[string]*TimRoomBean) []byte {
	t := &TimNodes{Ntype: ntype}
	if nodelist != nil {
		t.Nodelist = nodelist
	}
	if usermap != nil {
		t.Usermap = usermap
	}
	if roommap != nil {
		t.Roommap = roommap
	}
	buf := NewBuffer()
	buf.WriteByte(byte(TIMNODES))
	buf.Write(thrift.TEncode(t))
	return buf.Bytes()
}
