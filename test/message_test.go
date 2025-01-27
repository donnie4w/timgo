package test

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"
)

//tom1:UHuS8PoK2Mi
//tom2:QdH6CCms5FV

func TestMessageToUser(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.MessageToUser("QdH6CCms5FV", "hello123456哈", 0, 0, nil, nil)
	time.Sleep(2 * time.Second)
	t.Log(6)
}

func TestRevokeMessage(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.RevokeMessage(2, "QdH6CCms5FV", "", "", 0, 0)
	time.Sleep(3 * time.Second)
	t.Log(3)
}

func TestBurnMessage(t *testing.T) {
	tc := tclient()
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.BurnMessage(3, "UHuS8PoK2Mi", "", 0, 0)
	time.Sleep(3 * time.Second)
	t.Log(5)
}

func TestMessageToRoom(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.MessageToRoom("5pThTQ6Ycp1", fmt.Sprint("hello123456"), 0, 0, nil, nil)
	time.Sleep(3 * time.Second)
	t.Log(1)
}

func TestRevokeMessageRoom(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.RevokeMessage(1, "", "10001", "", 0, 0)
	time.Sleep(3 * time.Second)
	t.Log(4, nil)
}

func TestPullUserMessage(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.PullUserMessage("10001", 0, 10)
	time.Sleep(3 * time.Second)
	t.Log(1, nil)
}

func TestPullRoomMessage(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.PullRoomMessage("10001", 0, 10)
	time.Sleep(3 * time.Second)
	t.Log(2, nil)
}

func TestStreamToUser(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	bs := make([]byte, 1<<10)
	rand.Read(bs)
	tc.StreamToUser("QdH6CCms5FV", bs, 0, 0)
	time.Sleep(5 * time.Second)
	t.Log(1)
}
