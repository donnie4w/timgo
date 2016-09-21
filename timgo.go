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
	fmt.Println("versionName:", ProtocolversionName)

	flag.Parse()
	ip := flag.Arg(0)       //服务器ip地址
	port := flag.Arg(1)     //服务器端口
	fromname := flag.Arg(2) //登陆用户名
	toname := flag.Arg(3)   //发送信息对象
	tlsport := ""           //tls端口
	if flag.NArg() >= 5 {
		tlsport = flag.Arg(4)
	}

	for {
		clientTest(ip, port, fromname, "1234", toname, tlsport) /**1234为密码*/
	}

}

func clientTest(ip, port, fromname, pwd, toname, tlsport string) {
	conf := new(client.Conf)
	conf.Domain = "wuxiaodong"
	conf.Name = fromname
	conf.Pwd = pwd
	conf.Resource = "goclient"
	conf.AckListener = func(ab *TimAckBean) {} //ack监听
	conf.MessageListener = func(mbean *TimMBean) {
		fmt.Println("我收到[ "+mbean.GetFromTid().GetName()+" ]的消息：", mbean.GetBody())
	} //message 信息监听
	conf.PresenceListener = func(pbean *TimPBean) {
		fmt.Println(fmt.Sprint("我收到[ ", pbean.GetFromTid().GetName(), " ]的在线状态：", pbean.GetShow()))
	} //presence 好友在线状态监听
	conf.MessageListListener = func(mbeans []*TimMBean) {
		for _, mbean := range mbeans {
			fmt.Println("信息合流：我收到[ "+mbean.GetFromTid().GetName()+" ]的消息：", mbean.GetBody())
		}
	} // message 信息监听 列表
	conf.PresenceListListener = func(pbeans []*TimPBean) {
		for _, pbean := range pbeans {
			fmt.Println(fmt.Sprint("状态合流：我收到[ ", pbean.GetFromTid().GetName(), " ]的在线状态：", pbean.GetShow()))
		}
	} //presence 好友在线状态监听 列表

	tlsaddr := ""
	if tlsport != "" {
		tlsaddr = fmt.Sprint(ip, ":", tlsport)
	}

	/**
	*   如果tlsaddr!=""则采用tls传输
	*
	*
	 */
	cli, err := client.NewConn(fmt.Sprint(ip, ":", port), conf, tlsaddr)

	if err != nil {
		fmt.Println("create connect failed!", err)
		os.Exit(0)
	}

	for i := 1; i < 9999; i++ {
		time.Sleep(1000 * time.Millisecond)
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
