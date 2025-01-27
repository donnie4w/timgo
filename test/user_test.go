package test

import (
	"fmt"
	"github.com/donnie4w/timgo/stub"
	"testing"
	"time"
)

// 测试注册
func Test_register(t *testing.T) {
	register("tim1", "123")
	register("tim2", "123")
}

// 测试注册
func register(username, pwd string) {
	if ack, _ := tclient().Register(username, pwd, "tlnet.top"); ack != nil {
		fmt.Println(ack)
	}
}

func Test_login(t *testing.T) {
	//login("tim1", "123")
	login("tim2", "123")
	time.Sleep(10 * time.Minute)
}

// 登录
func login(username, pwd string) {
	tclient(50001).Login(username, pwd, "tlnet.top", "web", 1, nil)
}

func Test_LoginByToken(t *testing.T) {
	tclient().LoginByToken("aJsLRe81X1M", "webapp", 1, nil)
	time.Sleep(3 * time.Second)
}

func TestModify(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.ModifyPwd("123", "1234")
	time.Sleep(1 * time.Second)
	t.Log(5)
}

func TestNodeInfo(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	//tc.UserInfo("ijUgqi2oEs7", "f4X3mDaivzs")
	//tc.RoomInfo("fAyVgz7cpkw", "BaeoTNuaczC")
	name, nickname, brithday, gender, cover := "萌萌", "tomcat3", "2003-01-01", int8(0), "http://testwfs.tlnet.top:4660/1121372.jpeg"
	tc.ModifyUserInfo(&stub.TimUserBean{Name: &name, NickName: &nickname, Brithday: &brithday, Gender: &gender, Cover: &cover})
	time.Sleep(3 * time.Second)
	t.Log(1, nil)
}
