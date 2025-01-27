// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package test

import (
	"fmt"
	"github.com/donnie4w/timgo"
	. "github.com/donnie4w/timgo/stub"
	"log"
)

func tclient(ports ...int) *timgo.TimClient {
	port := 50001
	if len(ports) > 0 {
		port = ports[0]
	}
	return newTimClientWithHandle(false, "192.168.2.11", port)
}

func newTimClientWithHandle(tls bool, ip string, port int) *timgo.TimClient {
	tc := timgo.NewTimClient(tls, ip, port)
	var myAccount string
	//message消息
	tc.MessageHandler(func(tm *TimMessage) {
		if tm.MsType == 1 {
			log.Println("this is system message >>", "body>>>", tm.GetDataString())
		} else if tm.MsType == 2 {
			log.Println("this is user to user message")
		} else if tm.MsType == 3 {
			log.Println("this is room to user message")
		}
		if tm.MsType != 1 {
			switch tm.OdType {
			case 1: //常规消息
				log.Println("chat message user: from>>>", tm.FromTid.Node, ",to>>", tm.ToTid, ",room>>", tm.RoomTid)
				log.Println("chat message body>>>", *tm.DataString)
			case 2: //撤回消息
				log.Println("RevokeMessage>>>", *tm.Mid)
			case 3: //阅后即焚
				log.Println("BurnMessage>>>", *tm.Mid)
			case 4: //业务消息
				log.Println("business message>>>", tm)
			case 5: //流数据
				log.Println("stream message>>>", tm.DataBinary)
			default: //开发者自定义的消息
				log.Fatal("other message>>>", tm)
			}
		}
	})
	tc.PullmessageHandler(func(tm *TimMessageList) { log.Println("pull msg >>>", tm) })
	tc.OfflineMsgHandler(func(tmlist *TimMessageList) {
		fmt.Println("offline msg>>>", tmlist)
		for _, tm := range tmlist.MessageList {
			fmt.Println("offline msg ===> ", tm)
		}
	})

	//Indicates that the offline message is pushed
	//离线消息已经推送完毕
	tc.OfflineMsgEndHandler(func() {
		fmt.Println("OfflineMsgEndHandler")
		tc.Roster()                                            //Pull the roster 拉取花名册
		tc.UserRoom()                                          // pull the account of group 拉取群账号
		tc.BlockRoomList()                                     //blocklist of user group 用户群黑名单
		tc.BlockRosterList()                                   //blocklist of user 用户黑名单
		tc.BroadPresence(timgo.SUBSTATUS_REQ, 0, "I am busy😄") //when finish offline message recive,subcript and broad the presence
	})

	//记录状态订阅者
	submap := map[string]int8{}
	//状态消息
	tc.PresenceHandler(func(tp *TimPresence) {
		log.Println("PresenceHandler", tp.GetFromTid(), tp.GetToTid(), tp.GetStatus())
		if tp.SubStatus != nil {
			if tp.GetSubStatus() == timgo.SUBSTATUS_REQ {
				tc.PresenceToUser(tp.FromTid.Node, 0, "", timgo.SUBSTATUS_ACK, nil, nil)
			}
			if tp.GetSubStatus() == timgo.SUBSTATUS_REQ || tp.GetSubStatus() == timgo.SUBSTATUS_ACK {
				if _, ok := submap[tp.FromTid.Node]; ok {
					submap[tp.FromTid.Node] = 0
				}
			}
		}
		if tp.Offline != nil && tp.FromTid != nil && tp.FromTid.Node == myAccount {
			tc.BroadPresence(timgo.SUBSTATUS_REQ, 0, "I am busy😄")
		}
		log.Println("presence>>>", tp)
	})

	//现场流数据（即直播数据，实时语音视频等流数据）
	tc.StreamHandler(func(ts *TimStream) { log.Println("steamData>>>>", ts) })

	//ack message from server 服务反馈的信息
	tc.AckHandler(func(ta *TimAck) {
		log.Println("ack>>>", ta)
		switch timgo.TIMTYPE(ta.TimType) {
		case timgo.TIMMESSAGE:
			if !ta.Ok { //not ok 表示信息发送失败(注意，信息发送成功是没有ack的，服务器会推送发送用户的信息回来，信息会带上时间与id)
				log.Println("send message failed>>", *ta.Error.Code, ":", *ta.Error.Info)
			}
		case timgo.TIMPRESENCE:
			if !ta.Ok { //not ok，表示状态信息发送失败(注意，状态信息发送成功是没有ack的)
				log.Fatal("send presence failed>>", *ta.Error.Code, ":", *ta.Error.Info)
			}
		case timgo.TIMLOGOUT: // 强制下线
			log.Println("force to logout >>>", myAccount)
			tc.Logout() // 收到强制下线的指令后，主动退出登录
		case timgo.TIMAUTH:
			if ta.Ok { // 登录成功
				myAccount = ta.GetN()
				log.Println("login successful,my node is :", ta.GetN(), " ,uuid:", ta.GetT())
				tc.UserInfo(ta.GetN())
				tc.OfflineMsg() //when login successful, get the offline message 登录成功后，拉取离线信息
			} else {
				log.Fatal("login failed:", *ta.Error.Code, ":", *ta.Error.Info)
				tc.Logout()
				log.Fatal("login failed and logout")
			}
			// 虚拟房间操作回馈信息
		case timgo.TIMVROOM:
			if ta.Ok {
				switch ta.GetT() {
				case timgo.VRITURLROOM_REGISTER: //注册虚拟房间成功
					log.Println("register vriturl room ok>>", *ta.T, " >>", *ta.N)
				case timgo.VRITURLROOM_REMOVE: //删除虚拟房间成功
					log.Println("remove vriturl room ok>>", *ta.T, " >>", *ta.N)
				case timgo.VRITURLROOM_ADDAUTH: //加权成功
					log.Println("add auth to vriturl room ok>>", *ta.T, " >>", *ta.N)
				case timgo.VRITURLROOM_RMAUTH: //去权成功
					log.Println("cancel auth vriturl room process >>", *ta.T, " >>", *ta.N)
				case timgo.VRITURLROOM_SUB: //订阅成功
					log.Println("sub vriturl room process >>", *ta.T, " >>", *ta.N)
				case timgo.VRITURLROOM_SUBCANCEL: //取消订阅成功
					log.Println("sub cancel vriturl room process >>", *ta.T, " >>", *ta.N)
				default:
					log.Fatal("vriturl room process ok>>", *ta.T, " >>", *ta.N)
				}
			} else {
				log.Println("vriturl room process failed:", *ta.Error.Code, ":", *ta.Error.Info)
			}
		/*************************************************************/
		case timgo.TIMBUSINESS: //业务操作回馈
			if ta.Ok {
				log.Println("business process ok:", *ta.N)
				t := int32(*ta.T)
				switch t {
				case timgo.BUSINESS_REMOVEROSTER: //删除好友成功
					log.Println("romove friend successful:", *ta.N)
				case timgo.BUSINESS_BLOCKROSTER: //拉黑好友成功
					log.Println("block  successful:", *ta.N)
				case timgo.BUSINESS_NEWROOM: //新建群组成功
					log.Println("new group successful:", *ta.N)
				case timgo.BUSINESS_ADDROOM: //
					log.Println("join group successful:", *ta.N)
				case timgo.BUSINESS_PASSROOM: //申请加入群成功
					log.Println("new group successful:", *ta.N)
				case timgo.BUSINESS_NOPASSROOM: //申请加入群不成功
					log.Println("reject  group successful:", *ta.N)
				case timgo.BUSINESS_PULLROOM: //拉人入群
					log.Println("new group successful:", *ta.N)
				case timgo.BUSINESS_KICKROOM: //踢人出群
					log.Println("kick out of group successful:", *ta.N)
				case timgo.BUSINESS_BLOCKROOM:
					log.Println("block the group successful:", *ta.N)
				case timgo.BUSINESS_LEAVEROOM: //退群
					log.Println("leave group successful:", *ta.N)
				case timgo.BUSINESS_CANCELROOM: //注销群
					log.Println("cancel group successful:", *ta.N)
				default:
					log.Println("business successful >>", ta)
				}
			} else {
				log.Println("business process failed:", *ta.Error.Code, ":", *ta.Error.Info)
			}
		case timgo.TIMSTREAM:
			if !ta.Ok {
				log.Fatal("vriturl room process failed:", *ta.Error.Code, ":", *ta.Error.Info)
			}
		}
	})
	//批量账号信息返回
	tc.NodesHandler(func(tn *TimNodes) {
		switch tn.Ntype {
		case timgo.NODEINFO_ROSTER: //花名册返回
			log.Println("my roster >>>", tn)
			if tn.Nodelist != nil {
				tc.UserInfo(tn.Nodelist...) //获取用户详细资料
			}
		case timgo.NODEINFO_ROOM: //用户的群账号返回
			log.Println("my groups >>>", tn)
			if tn.Nodelist != nil {
				tc.RoomInfo(tn.Nodelist...) //获取群详细资料
				for _, node := range tn.Nodelist {
					tc.RoomUsers(node) //获取群的成员
				}
			}
		case timgo.NODEINFO_ROOMMEMBER: //群成员账号返回
			log.Println("group member ack >>>", tn)
		case timgo.NODEINFO_USERINFO: //用户信息返回
			log.Println("userinfo ack >>>")
			for k, v := range tn.Usermap {
				log.Println(k, ">>", v.GetName(), " ", v.GetNickName(), " ", v.GetBrithday(), " ", v.GetGender(), " ", v.GetCover(), " ", v.GetArea(), " ", v.GetPhotoTidAlbum())
			}
		case timgo.NODEINFO_ROOMINFO: //群信息返回
			log.Println("groupinfo ack >>>")
			for k, v := range tn.Roommap {
				log.Println(k, " >>", v.GetTopic(), " ", v.GetFounder(), " ", v.GetManagers(), " ", v.GetCreatetime(), " ", v.GetLabel())
			}
		case timgo.NODEINFO_BLOCKROSTERLIST: //用户黑名单
			log.Println("block roster list ack >>>", tn.Nodelist)
		case timgo.NODEINFO_BLOCKROOMLIST: //用户拉黑群的群账号
			log.Println("block room list ack >>>", tn.Nodelist)
		case timgo.NODEINFO_BLOCKROOMMEMBERLIST: //群拉黑账号名单
			log.Println("block room member list ack >>>", tn.Nodelist)
		}
	})
	return tc
}
