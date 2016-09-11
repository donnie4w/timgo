package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"timgo.client"
	. "timgo.protocol"
)

func main() {
	flag.Parse()
	ip := flag.Arg(0)
	port := flag.Arg(1)
	fromname := flag.Arg(2)
	toname := flag.Arg(3)
	for {
		clientTest(ip, port, fromname, "1234", toname)
	}

}

func clientTest(ip, port, fromname, pwd, toname string) {
	conf := new(client.Conf)
	conf.Domain = "wuxiaodong"
	conf.Name = fromname
	conf.Pwd = pwd
	conf.Resource = "goclient"
	conf.AckListener = func(ab *TimAckBean) {}
	conf.MessageListener = func(mbean *TimMBean) {
		fmt.Println("我收到[ "+mbean.GetFromTid().GetName()+" ]的消息：", mbean.GetBody())
	}
	conf.PresenceListener = func(pbean *TimPBean) {
		fmt.Println(fmt.Sprint("我收到[ ", pbean.GetFromTid().GetName(), " ]的在线状态：", pbean.GetShow()))
	}
	conf.MessageListListener = func(mbeans []*TimMBean) {
		for _, mbean := range mbeans {
			fmt.Println("信息合流：我收到[ "+mbean.GetFromTid().GetName()+" ]的消息：", mbean.GetBody())
		}
	}
	conf.PresenceListListener = func(pbeans []*TimPBean) {
		for _, pbean := range pbeans {
			fmt.Println(fmt.Sprint("状态合流：我收到[ ", pbean.GetFromTid().GetName(), " ]的在线状态：", pbean.GetShow()))
		}
	}

	cli, err := client.NewConn(fmt.Sprint(ip, ":", port), conf)

	if err != nil {
		fmt.Println("create connect failed!", err)
		os.Exit(0)
	}

	for i := 1; i < 9999; i++ {
		time.Sleep(100 * time.Millisecond)
		msg := fmt.Sprint(fromname, "=====>", toname, ":", i)
		err = cli.Sendmsg(toname, &msg)
		if err != nil {
			goto END_
		}
		show := fmt.Sprint("我在打电话:", i)
		err = cli.SendPresence(toname, &show)
		if err != nil {
			goto END_
		}
	}
END_:
	time.Sleep(5 * time.Second)
	return
}
