// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package timgo

import (
	"flag"
	"fmt"
	"time"

	"github.com/donnie4w/simplelog/logging"
	. "github.com/donnie4w/timgo/stub"
)

func newTimClientWithHandle(ip string, port int, tls bool) *TimClient {
	tc := NewTimClient(ip, port, tls)
	var myAccount string
	//message消息
	tc.MessageHandler(func(tm *TimMessage) {
		if tm.MsType == 1 {
			logging.Debug("this is system message >>", "body>>>", tm)
		} else if tm.MsType == 2 {
			logging.Debug("this is user to user message")
		} else if tm.MsType == 3 {
			logging.Debug("this is room to user message")
		}
		if tm.MsType != 1 {
			switch tm.OdType {
			case 1: //常规消息
				logging.Debug("chat message user: from>>>", tm.FromTid.Node, ",to>>", tm.ToTid, ",room>>", tm.RoomTid)
				logging.Debug("chat message body>>>", *tm.DataString)
			case 2: //撤回消息
				logging.Debug("RevokeMessage>>>", *tm.Mid)
			case 3: //阅后即焚
				logging.Debug("BurnMessage>>>", *tm.Mid)
			case 4: //业务消息
				logging.Debug("business message>>>", tm)
			case 5: //流数据
				logging.Debug("stream message>>>", tm.DataBinary)
			default: //开发者自定义的消息
				logging.Fatal("other message>>>", tm)
			}
		}
	})
	tc.PullmessageHandler(func(tm *TimMessageList) { logging.Debug("pull msg >>>", tm) })
	tc.OfflineMsgHandler(func(tm *TimMessageList) {
		logging.Debug("offline msg>>>", tm)
		for _, s := range tm.MessageList {
			logging.Debug("offline msg ===> ", *s.DataString)
		}
	})

	//Indicates that the offline message is pushed
	//离线消息已经推送完毕
	tc.OfflineMsgEndHandler(func() {
		logging.Debug("OfflineMsgEndHandler")
		tc.Roster()                                      //Pull the roster 拉取花名册
		tc.UserRoom()                                    // pull the account of group 拉取群账号
		tc.BlockRoomList()                               //blocklist of user group 用户群黑名单
		tc.BlockRosterList()                             //blocklist of user 用户黑名单
		tc.BroadPresence(SUBSTATUS_REQ, 0, "I am busy😄") //when finish offline message recive,subcript and broad the presence
	})

	//记录状态订阅者
	submap := map[string]int8{}
	//状态消息
	tc.PresenceHandler(func(tp *TimPresence) {
		logging.Debug(tp.FromTid.Node, " presnce")
		if tp.SubStatus != nil {
			if *tp.SubStatus == SUBSTATUS_REQ {
				tc.PresenceToUser(tp.FromTid.Node, 0, "", SUBSTATUS_ACK, nil, nil)
			}
			if *tp.SubStatus == SUBSTATUS_REQ || *tp.SubStatus == SUBSTATUS_ACK {
				if _, ok := submap[tp.FromTid.Node]; ok {
					submap[tp.FromTid.Node] = 0
				}
			}
		}
		if tp.Offline != nil && tp.FromTid.Node == myAccount {
			tc.BroadPresence(SUBSTATUS_REQ, 0, "I am busy😄")
		}
		logging.Debug("presence>>>", tp)
	})

	//现场流数据（即直播数据，实时语音视频等流数据）
	tc.StreamHandler(func(ts *TimStream) { logging.Debug("steamData>>>>", ts) })

	//ack message from server 服务反馈的信息
	tc.AckHandler(func(ta *TimAck) {
		logging.Debug("ack>>>", ta)
		switch TIMTYPE(ta.TimType) {
		case TIMMESSAGE:
			if !ta.Ok { //not ok 表示信息发送失败(注意，信息发送成功是没有ack的，服务器会推送发送用户的信息回来，信息会带上时间与id)
				logging.Error("send message failed>>", *ta.Error.Code, ":", *ta.Error.Info)
			}
		case TIMPRESENCE:
			if !ta.Ok { //not ok，表示状态信息发送失败(注意，状态信息发送成功是没有ack的)
				logging.Error("send presence failed>>", *ta.Error.Code, ":", *ta.Error.Info)
			}
		case TIMLOGOUT: // 强制下线
			logging.Debug("force to logout >>>", myAccount)
			tc.Logout() // 收到强制下线的指令后，主动退出登录
		case TIMAUTH:
			if ta.Ok { // 登录成功
				myAccount = *ta.N
				logging.Debug("login successful,my node is :", myAccount)
				tc.UserInfo(myAccount)
				tc.OfflineMsg() //when login successful, get the offline message 登录成功后，拉取离线信息
			} else {
				logging.Error("login failed:", *ta.Error.Code, ":", *ta.Error.Info)
				tc.Logout()
				logging.Error("login failed and logout")
			}
			// 虚拟房间操作回馈信息
		case TIMVROOM:
			if ta.Ok {
				switch *ta.T {
				case VRITURLROOM_REGISTER: //注册虚拟房间成功
					logging.Info("register vriturl room ok>>", *ta.T, " >>", *ta.N)
				case VRITURLROOM_REMOVE: //删除虚拟房间成功
					logging.Info("remove vriturl room ok>>", *ta.T, " >>", *ta.N)
				case VRITURLROOM_ADDAUTH: //加权成功
					logging.Info("add auth to vriturl room ok>>", *ta.T, " >>", *ta.N)
				case VRITURLROOM_RMAUTH: //去权成功
					logging.Info("cancel auth vriturl room process >>", *ta.T, " >>", *ta.N)
				case VRITURLROOM_SUB: //订阅成功
					logging.Info("sub vriturl room process >>", *ta.T, " >>", *ta.N)
				case VRITURLROOM_SUBCANCEL: //取消订阅成功
					logging.Info("sub cancel vriturl room process >>", *ta.T, " >>", *ta.N)
				default:
					logging.Fatal("vriturl room process ok>>", *ta.T, " >>", *ta.N)
				}
			} else {
				logging.Info("vriturl room process failed:", *ta.Error.Code, ":", *ta.Error.Info)
			}
		/*************************************************************/
		case TIMBUSINESS: //业务操作回馈
			if ta.Ok {
				logging.Info("business process ok:", *ta.N)
				t := int32(*ta.T)
				switch t {
				case BUSINESS_REMOVEROSTER: //删除好友成功
					logging.Debug("romove friend successful:", *ta.N)
				case BUSINESS_BLOCKROSTER: //拉黑好友成功
					logging.Debug("block  successful:", *ta.N)
				case BUSINESS_NEWROOM: //新建群组成功
					logging.Debug("new group successful:", *ta.N)
				case BUSINESS_ADDROOM: //
					logging.Debug("join group successful:", *ta.N)
				case BUSINESS_PASSROOM: //申请加入群成功
					logging.Debug("new group successful:", *ta.N)
				case BUSINESS_NOPASSROOM: //申请加入群不成功
					logging.Debug("reject  group successful:", *ta.N)
				case BUSINESS_PULLROOM: //拉人入群
					logging.Debug("new group successful:", *ta.N)
				case BUSINESS_KICKROOM: //踢人出群
					logging.Debug("kick out of group successful:", *ta.N)
				case BUSINESS_BLOCKROOM:
					logging.Debug("block the group successful:", *ta.N)
				case BUSINESS_LEAVEROOM: //退群
					logging.Debug("leave group successful:", *ta.N)
				case BUSINESS_CANCELROOM: //注销群
					logging.Debug("cancel group successful:", *ta.N)
				default:
					logging.Debug("business successful >>", ta)
				}
			} else {
				logging.Info("business process failed:", *ta.Error.Code, ":", *ta.Error.Info)
			}
		case TIMSTREAM:
			if !ta.Ok {
				logging.Error("vriturl room process failed:", *ta.Error.Code, ":", *ta.Error.Info)
			}
		}
	})
	//批量账号信息返回
	tc.NodesHandler(func(tn *TimNodes) {
		switch tn.Ntype {
		case NODEINFO_ROSTER: //花名册返回
			logging.Debug("my roster >>>", tn)
			if tn.Nodelist != nil {
				tc.UserInfo(tn.Nodelist...) //获取用户详细资料
			}
		case NODEINFO_ROOM: //用户的群账号返回
			logging.Debug("my groups >>>", tn)
			if tn.Nodelist != nil {
				tc.RoomInfo(tn.Nodelist...) //获取群详细资料
				for _, node := range tn.Nodelist {
					tc.RoomUsers(node) //获取群的成员
				}
			}
		case NODEINFO_ROOMMEMBER: //群成员账号返回
			logging.Debug("group member ack >>>", tn)
		case NODEINFO_USERINFO: //用户信息返回
			logging.Debug("userinfo ack >>>")
			for k, v := range tn.Usermap {
				logging.Debug(k, ">>", v.GetName(), " ", v.GetNickName(), " ", v.GetBrithday(), " ", v.GetGender(), " ", v.GetCover(), " ", v.GetArea(), " ", v.GetPhotoTidAlbum())
			}
		case NODEINFO_ROOMINFO: //群信息返回
			logging.Debug("groupinfo ack >>>")
			for k, v := range tn.Roommap {
				logging.Debug(k, " >>", v.GetTopic(), " ", v.GetFounder(), " ", v.GetManagers(), " ", v.GetCreatetime(), " ", v.GetLabel())
			}
		case NODEINFO_BLOCKROSTERLIST: //用户黑名单
			logging.Debug("block roster list ack >>>", tn.Nodelist)
		case NODEINFO_BLOCKROOMLIST: //用户拉黑群的群账号
			logging.Debug("block room list ack >>>", tn.Nodelist)
		case NODEINFO_BLOCKROOMMEMBERLIST: //群拉黑账号名单
			logging.Debug("block room member list ack >>>", tn.Nodelist)
		}
	})
	return tc
}

func main() {
	// subVirtual()
	name := flag.String("n", "", "")
	pwd := flag.String("pwd", "", "")
	p := flag.Int("p", 0, "")
	flag.Parse()
	tc := newTimClientWithHandle("tim.tlnet.top", *p, true)
	tc.Login(*name, *pwd, "tlnet.top", "android", 1, nil)
	select {}
}

func loginone() *TimClient {
	name := flag.String("n", "", "")
	pwd := flag.String("pwd", "", "")
	p := flag.Int("p", 0, "")
	flag.Parse()
	tc := newTimClientWithHandle("192.168.2.11", *p, false)
	tc.Login(*name, *pwd, "tlnet.top", "android", 1, nil)
	return tc
}

func subVirtual() {
	vr := flag.String("v", "", "")
	name := flag.String("n", "", "")
	pwd := flag.String("pwd", "", "")
	p := flag.Int("p", 0, "")
	flag.Parse()
	tc := newTimClientWithHandle("192.168.2.11", *p, false)
	tc.Login(*name, *pwd, "tlnet.top", "android", 1, nil)
	<-time.After(time.Second)
	tc.VirtualroomSub(*vr)
}

func loginMutli() {
	f := flag.Int("f", 0, "")
	t := flag.Int("t", 0, "")
	ip := flag.String("ip", "", "")
	p := flag.Int("p", 0, "")
	flag.Parse()
	if f == nil || t == nil || ip == nil || p == nil {
		logging.Error(">>>", f, ",", t, ",", ip, ",", p)
		panic("error params")
	}
	logging.Debug("loginMutli: f>>", *f, ",t>>", *t, ",ip>>", *ip, ",port>>", *p)
	for i := *f; i <= *t; i++ {
		<-time.After(100 * time.Millisecond)
		go newTimClientWithHandle(*ip, *p, false).Login(fmt.Sprint("tim", i), "123", "tlnet.top", "android", 1, nil)
	}
}

func newAccount() {
	f := flag.Int("f", 0, "")
	t := flag.Int("t", 0, "")
	ip := flag.String("ip", "", "")
	p := flag.Int("p", 0, "")
	flag.Parse()
	if f == nil || t == nil || ip == nil || p == nil {
		logging.Error(">>>", f, ",", t, ",", ip, ",", p)
		panic("error params")
	}
	logging.Debug("newAccount: f>>", *f, ",t>>", *t, ",ip>>", ",port>>", *p)
	registermulti(*f, *t, *ip, *p)
}

func registermulti(f, t int, ip string, p int) {
	for i := f; i <= t; i++ {
		<-time.After(5 * time.Millisecond)
		go func(i int) {
			sc := newTimClientWithHandle(ip, p, false)
			if ack, err := sc.Register(fmt.Sprint("tim", i), "123", "tlnet.top"); err == nil && ack != nil && ack.Ok {
				logging.Debug("register successful")
				logging.Debug("node>>>", *ack.N)
			} else if ack != nil {
				logging.Debug(fmt.Sprint("tim", i), ",register failed>>", *ack.Error.Info)
			}
			sc.Login(fmt.Sprint("tim", i), "123", "tlnet.top", "android", 1, nil)
		}(i)
	}
}
