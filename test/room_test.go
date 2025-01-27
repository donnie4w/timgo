package test

import (
	"github.com/donnie4w/timgo"
	"testing"
	"time"
)

// 5pThTQ6Ycp1
func TestNewRoom(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.NewRoom(timgo.ROOM_OPEN, "tim group")
	time.Sleep(3 * time.Second)
	t.Log("TestNewRoom")
}

// private room : AaAGFd4JaHf
func TestAddRoom(t *testing.T) {
	tc := tclient()
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.AddRoom("5pThTQ6Ycp1", "i am tim")
	time.Sleep(3 * time.Second)
	t.Log("TestAddRoom")
}

func TestPullInRoom(t *testing.T) {
	tc := tclient()
	tc.Login("test1", "123", "tlnet.top", "android", 1, nil)
	tc.PullInRoom("AaAGFd4JaHf", "UHuS8PoK2Mi")
	time.Sleep(3 * time.Second)
	t.Log("TestPullInRoom")
}

func TestBlockRoom(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.BlockRoom("AaAGFd4JaHf")
	time.Sleep(3 * time.Second)
	t.Log("TestBlockRoom")
}

func TestLeaveRoom(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.LeaveRoom("AaAGFd4JaHf")
	time.Sleep(3 * time.Second)
	t.Log("TestLeaveRoom")
}

func TestBlockRoomMember(t *testing.T) {
	tc := tclient()
	tc.Login("test1", "123", "tlnet.top", "android", 1, nil)
	tc.BlockRoomMember("AaAGFd4JaHf", "QdH6CCms5FV")
	time.Sleep(3 * time.Second)
	t.Log("TestBlockRoomMember")
}

func TestKickRoom(t *testing.T) {
	tc := tclient()
	tc.Login("test1", "123", "tlnet.top", "android", 1, nil)
	tc.KickRoom("AaAGFd4JaHf", "QdH6CCms5FV")
	time.Sleep(3 * time.Second)
	t.Log("TestKickRoom")
}

func TestRosterPullInRoom(t *testing.T) {
	tc := tclient()
	tc.Login("test1", "123", "tlnet.top", "android", 1, nil)
	tc.PullInRoom("AaAGFd4JaHf", "UHuS8PoK2Mi")
	time.Sleep(3 * time.Second)
	t.Log("TestRosterPullInRoom")
}

func TestBlockRoomMemberlist(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.BlockRoomMemberlist("AaAGFd4JaHf")
	time.Sleep(3 * time.Second)
	t.Log("TestBlockRoomMemberlist")
}
