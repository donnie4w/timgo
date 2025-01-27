package test

import (
	"testing"
	"time"
)

//tim1 UHuS8PoK2M9
//tim2 QdH6CCms5Ex

func TestAddroster(t *testing.T) {
	tc := tclient()
	tc.Login("tim2", "123", "tlnet.top", "android", 1, nil)
	tc.Addroster("UHuS8PoK2M9", "123")
	time.Sleep(3 * time.Second)
	t.Log(4)
}

func TestRmroster(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.Rmroster("QdH6CCms5Ex")
	time.Sleep(3 * time.Second)
	t.Log(3)
}

func TestBlockroster(t *testing.T) {
	tc := tclient()
	tc.Login("tim1", "123", "tlnet.top", "android", 1, nil)
	tc.Blockroster("QdH6CCms5Ex")
	time.Sleep(3 * time.Second)
	t.Log(3)
}
