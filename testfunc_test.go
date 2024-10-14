// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package timgo

import (
	"fmt"
	"testing"
	"time"

	"github.com/donnie4w/gofer/util"
	"github.com/donnie4w/simplelog/logging"
)

func newaccount(number int, preusername string) {
	for i := 0; i < number; i++ {
		<-time.After(10 * time.Millisecond)
		go func(i int) {
			tc := newTimClientWithHandle("192.168.2.11", 5080, false)
			account := fmt.Sprint(preusername, i)
			if ack, err := tc.Register(account, "123", "tlnet.top"); err == nil && ack != nil && ack.Ok {
				logging.Debug("register successful")
				logging.Debug("node>>>", *ack.N)
			} else if ack != nil {
				logging.Debug("register failed>>", *ack.Error.Info)
			}
			tc.Login(account, "123", "tlnet.top", "android", 1, nil)
			<-time.After(time.Second)
			tc.AddRoom("JXjh2vcocNk", "I am abc")
		}(i)
	}
}

func newlogin(from, to int, port int, name string) {
	for i := from; i < to; i++ {
		<-time.After(10 * time.Millisecond)
		go func(i int) {
			tc := newTimClientWithHandle("192.168.2.11", 5080, false)
			account := fmt.Sprint(name, i)
			tc.Login(account, "123", "tlnet.top", "android", 1, nil)
			<-time.After(time.Second)
			// tc.AddROOM("BaeoTNuaczC")
		}(i)
	}
}

func TestLogin10000(t *testing.T) {
	newlogin(5000, 10000, 5082, "new")
	select {}
}

func TestLogin5000(t *testing.T) {
	newlogin(0, 5000, 5080, "new")
	select {}
}

func TestLogin1000(t *testing.T) {
	newlogin(0, 1, 5080, "qian")
	select {}
}

func TestNewaccount(t *testing.T) {
	newaccount(100, "qian")
	select {}
}

func Test_Register(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	if ack, _ := tc.Register("tom1", "123", "tlnet.top"); ack != nil {
		fmt.Println(ack)
	}
}

func Test_Login(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tom2", "123", "tlnet.top", "web", 1, nil)
	<-time.After(300 * time.Second)
}

func TestLogin(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5120, false)
	/*————————————————————————————————————————————————————————————————————————————————————————*/
	// if ack, err := tc.Register("test5", "123", "tlnet.top"); err == nil && ack != nil && ack.Ok {
	// 	logging.Debug("register successful")
	// 	logging.Debug("node>>>", *ack.N)
	// 	tc.Login("test4", "123", "tlnet.top", "android", 1,nil)
	// } else if ack != nil {
	// 	logging.Debug("register failed>>", *ack.Error.Info)
	// }
	tc.Login("test5", "123", "tlnet.top", "android", 1, nil)
	// tc.NewROOM(2, "open2")
	// tc.Addroster("ijUgqi2oEs7", "i am test1")
	// tc.MessageToUser("ijUgqi2oEs7", fmt.Sprint("hello123456"), 0, 0, nil, nil)
	// token, _ := tc.Token("test1", "123", "tlnet.top")
	// tc.LoginByToken(token, "android") //login by token
	// for i := 1; i < 10; i++ {
	// 	<-time.After(1 * time.Second)
	// 	tc.MessageToRoom("BaeoTNuaczC", fmt.Sprint("hello123456",i), 0, 0, nil, nil)
	// }
	// tc.StreamToUser("ijUgqi2oEs7", []byte("this is stream package"))
	// tc.StreamToRoom("fAyVgz7cpkw", []byte{1, 2})
	// tc.RevokeMessage(8, "d4RCZD2bxiW", "", 0, 0, nil, nil) // 撤回
	// tc.BurnMessage(9, "ijUgqi2oEs7", "", 0, 0, nil, nil) //阅后即焚
	// tc.Rmroster("ijUgqi2oEs7") //删除好友
	// tc.Blockroster("ijUgqi2oEs7") //拉入黑名单
	// tc.AddRoom("fAyVgz7cpkw", "i am tom") //加入群
	// tc.PullInRoom("fAyVgz7cpkw", "ijUgqi2oEs7")

	// tc.ROOMusers("BaeoTNuaczC")
	// for i := 0; i < 1000; i++ {
	// 	<-time.After(1 * time.Millisecond)
	// 	tc.MessageToRoom("BaeoTNuaczC", "hello123456", 0, 0, nil, nil)
	// }
	// tc.RevokeMessage(3, "", "fAyVgz7cpkw", "RevokeMessage", 0, 0)
	// tc.AddROOM("BaeoTNuaczC")
	// tc.LeaveRoom("fAyVgz7cpkw")
	// tc.PullRoomMessage("JXjh2vcocNk", 0, 4)
	<-time.After(1 * time.Second)
	logging.Debug(1, nil)
}

func TestRegistermulti(t *testing.T) {
	registermulti(1, 100, "192.168.2.11", 5080)
	<-time.After(1 * time.Second)
	logging.Debug(1, nil)
}

//test1 d4RCZD2bxiW
//test2 ijUgqi2oEs7
//test3 f4X3mDaivzs
//test4 UiBHQN1RLUj
//test5 Z4FCe9cZRFP

//private fAyVgz7cpkw
//openroom BaeoTNuaczC  10000
//openroom JXjh2vcocNk  1000

func TestVirtualroomRegister(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5081, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.VirtualroomRegister()
	<-time.After(2 * time.Second)
	logging.Debug(5)
}

func TestPushStream(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5081, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	for i := 0; i < 1<<1; i++ {
		<-time.After(100 * time.Microsecond)
		tc.PushStream("4cV9YaHc1sd", []byte{uint8(i), uint8(i) + 1, uint8(i) + 2, uint8(i) + 3, uint8(i) + 4, uint8(i) + 5}, 1)
	}
	<-time.After(10 * time.Second)
	logging.Debug(3)
}

func TestVirtualroomAddAuth(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.VirtualroomSub("AoCr4JXU9KM")
	<-time.After(5 * time.Second)
	tc.VirtualroomAddAuth("4t5b1BsZqV2", "ijUgqi2oEs7")
	tc.VirtualroomDelAuth("JaCNH7vfTLn", "f4X3mDaivzs")
	<-time.After(2 * time.Second)
	logging.Debug(1, nil)
}

func TestModify(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.ModifyPwd("123", "123")
	<-time.After(1 * time.Second)
	logging.Debug(5)
}

func TestVriturlSub(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.VirtualroomSub("jBfuhomE1um")
	// tc.VirtualroomSubCancel("AoCr4JXU9KM")
	<-time.After(2 * time.Second)
	logging.Debug(1, nil)
}

func TestVirtualroomRemove(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.VirtualroomRemove("Q6uWRj6inkY")
	<-time.After(3 * time.Second)
	logging.Debug(1, nil)
}

func TestNodeInfo(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("test1", "123", "tlnet.top", "android", 1, nil)
	tc.UserInfo("ijUgqi2oEs7", "f4X3mDaivzs")
	tc.RoomInfo("fAyVgz7cpkw", "BaeoTNuaczC")
	// name, nickname, brithday, gender, cover := "tom3", "tomcat3", "2003-01-01", int8(0), "https://xxx"
	// tc.ModifyUserInfo(&TimUserBean{Name: &name, NickName: &nickname, Brithday: &brithday, Gender: &gender, Cover: &cover})
	<-time.After(3 * time.Second)
	logging.Debug(1, nil)
}

func TestPullUserMessage(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5000, false)
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.PullUserMessage("10001", 0, 10)
	<-time.After(3 * time.Second)
	logging.Debug(1, nil)
}

func TestPullRoomMessage(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5000, false)
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.PullRoomMessage("10000000001", 0, 10)
	<-time.After(3 * time.Second)
	logging.Debug(1, nil)
}

func TestRevokeMessage(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5000, false)
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.RevokeMessage(1, "10001", "", "", 0, 0)
	<-time.After(3 * time.Second)
	logging.Debug(2)
}

func TestRevokeMessageRoom(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5000, false)
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.RevokeMessage(7, "", "10000000001", "", 0, 0)
	<-time.After(3 * time.Second)
	logging.Debug(1, nil)
}

func TestBurnMessage(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5082, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.BurnMessage(4, "10002", "", 0, 0)
	<-time.After(3 * time.Second)
	logging.Debug(2)
}

func TestMessageToUser(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5081, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.MessageToUser("QdH6CCms5FV", fmt.Sprint("hello123456哈"), 0, 0, nil, nil)
	<-time.After(3 * time.Second)
	logging.Debug(2)
}

func TestMessageToRoom(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5000, false)
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.MessageToRoom("10000000001", fmt.Sprint("hello123456"), 0, 0, nil, nil)
	<-time.After(3 * time.Second)
	logging.Debug(1)
}

func TestStreamToUser(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 443, true)
	tc.Login("tim3", "123", "tlnet.top", "android", 1, nil)
	bs, _ := util.ReadFile(`D:\workspace\donnie4w_go\src\github.com\webtim\webtim_linux`)
	tc.StreamToUser("UHuS8PoK2Mi", bs[:], 0, 0)
	<-time.After(2 * time.Second)
	logging.Debug(1)
}

// tim1 UHuS8PoK2Mi
// tim2 QdH6CCms5FV

func TestAddroster(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim8", "123", "tlnet.top", "android", 1, nil)
	tc.Addroster("QdH6CCms5FV", "123")
	<-time.After(3 * time.Second)
	logging.Debug(4)
}

func TestRmroster(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.Rmroster("QdH6CCms5FV")
	<-time.After(3 * time.Second)
	logging.Debug(3)
}

func TestBlockroster(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.Blockroster("QdH6CCms5FV")
	<-time.After(3 * time.Second)
	logging.Debug(3)
}

func TestRosterPullInRoom(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("test1", "123", "tlnet.top", "android", 1, nil)
	tc.PullInRoom("AaAGFd4JaHf", "UHuS8PoK2Mi")
	<-time.After(3 * time.Second)
	logging.Debug(1, nil)
}

func TestRoomNewRoom(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("qian1", "123", "tlnet.top", "android", 1, nil)
	tc.NewRoom(ROOM_OPEN, "tim group")
	<-time.After(3 * time.Second)
	logging.Debug(4)
}

// private room : AaAGFd4JaHf
func TestRoomAddRoom(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.AddRoom("AaAGFd4JaHf", "i am tim1")
	<-time.After(3 * time.Second)
	logging.Debug(1, nil)
}

func TestRoomPullInRoom(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("test1", "123", "tlnet.top", "android", 1, nil)
	tc.PullInRoom("AaAGFd4JaHf", "UHuS8PoK2Mi")
	<-time.After(3 * time.Second)
	logging.Debug(1, nil)
}

func TestBlockRoom(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.BlockRoom("AaAGFd4JaHf")
	<-time.After(3 * time.Second)
	logging.Debug(4)
}

func TestLeaveRoom(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.LeaveRoom("AaAGFd4JaHf")
	<-time.After(3 * time.Second)
	logging.Debug(4)
}

func TestBlockRoomMember(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("test1", "123", "tlnet.top", "android", 1, nil)
	tc.BlockRoomMember("AaAGFd4JaHf", "QdH6CCms5FV")
	<-time.After(3 * time.Second)
	logging.Debug(1, nil)
}

func TestKickRoom(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("test1", "123", "tlnet.top", "android", 1, nil)
	tc.KickRoom("AaAGFd4JaHf", "QdH6CCms5FV")
	<-time.After(3 * time.Second)
	logging.Debug(1, nil)
}

func TestBlockRoomMemberlist(t *testing.T) {
	tc := newTimClientWithHandle("192.168.2.11", 5080, false)
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.BlockRoomMemberlist("AaAGFd4JaHf")
	<-time.After(3 * time.Second)
	logging.Debug(2)
}
