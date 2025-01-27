package test

import (
	"testing"
	"time"
)

// 虚拟房间注册
func TestVirtualroomRegister(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.VirtualroomRegister()
	time.Sleep(2 * time.Second)
}

func TestVriturlSub(t *testing.T) {
	tc := tclient()
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.VirtualroomSub("5UuHvyshZND")
	//tc.VirtualroomSubCancel("AoCr4JXU9KM")
	time.Sleep(5 * time.Minute)
	t.Log(1, nil)
}

func TestVirtualroomRemove(t *testing.T) {
	tc := tclient()
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.VirtualroomRemove("Q6uWRj6inkY")
	time.Sleep(3 * time.Second)
	t.Log(1, nil)
}

// 虚拟房间推流
func TestPushStream(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Microsecond)
		tc.PushStream("KonwQpiwisv", []byte{uint8(i), uint8(i) + 1, uint8(i) + 2, uint8(i) + 3, uint8(i) + 4, uint8(i) + 5}, 1)
	}
	time.Sleep(10 * time.Second)
	t.Log(3)
}

func Test_BigDataBinary(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	err := tc.BigDataBinary("ZsV5WMTqYAv", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	time.Sleep(3 * time.Second)
	t.Log(1, err)
}

func Test_BigDataString(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.BigDataString("ZsV5WMTqYAv", "123456789")
	time.Sleep(3 * time.Second)
	t.Log(1, nil)
}
