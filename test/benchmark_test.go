package test

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func registermulti(f, t int, ip string, p int) {
	for i := f; i <= t; i++ {
		time.Sleep(5 * time.Millisecond)
		go func(i int) {
			tc := newTimClientWithHandle(false, ip, p)
			if ack, err := tc.Register(fmt.Sprint("tim", i), "123", "tlnet.top"); err == nil && ack != nil && ack.Ok {
				log.Println("register successful")
				log.Println("node>>>", *ack.N)
			} else if ack != nil {
				log.Println(fmt.Sprint("tim", i), ",register failed>>", *ack.Error.Info)
			}
			tc.Login(fmt.Sprint("tim", i), "123", "tlnet.top", "android", 1, nil)
		}(i)
	}
}

func Test_Registermulti(t *testing.T) {
	registermulti(1, 100, "192.168.2.11", 5080)
	time.Sleep(1 * time.Second)
	t.Log(1, nil)
}

func newlogin(from, to int, port int, name string) {
	for i := from; i < to; i++ {
		time.Sleep(10 * time.Millisecond)
		go func(i int) {
			tc := newTimClientWithHandle(false, "192.168.2.11", port)
			account := fmt.Sprint(name, i)
			tc.Login(account, "123", "tlnet.top", "android", 1, nil)
			time.Sleep(time.Second)
		}(i)
	}
}

func newaccount(number int, preusername string) {
	for i := 0; i < number; i++ {
		time.Sleep(10 * time.Millisecond)
		go func(i int) {
			tc := tclient()
			account := fmt.Sprint(preusername, i)
			if ack, err := tc.Register(account, "123", "tlnet.top"); err == nil && ack != nil && ack.Ok {
				fmt.Println("register successful")
				fmt.Println("node>>>", *ack.N)
			} else if ack != nil {
				fmt.Println("register failed>>", *ack.Error.Info)
			}
			tc.Login(account, "123", "tlnet.top", "android", 1, nil)
			time.Sleep(time.Second)
			tc.AddRoom("JXjh2vcocNk", "I am abc")
		}(i)
	}
}

func TestLogin10000(t *testing.T) {
	newlogin(5000, 10000, 5082, "new")
	select {}
}
