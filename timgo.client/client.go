package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/go-logger/logger"
	. "timgo.protocol"
)

type FLOW string

const (
	START  FLOW = "start"
	AUTH   FLOW = "auth"
	NOAUTH FLOW = "noauth"
	CLOSE  FLOW = "close"

	CONNECT_START FLOW = "connect_start"
	CONNECT_RUN   FLOW = "connect_run"
	CONNECT_STOP  FLOW = "connect_stop"
)

type Connect struct {
	Client      *ITimClient
	FlowConnect FLOW
	Super       *Cli
}

func (this *Connect) Close() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("sendmsg,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if this.Client != nil && this.Client.Transport != nil && this.FlowConnect != CONNECT_STOP {
		this.FlowConnect = CONNECT_STOP
		this.Client.Transport.Close()
	}
}

func (this *Connect) setITimClient(client *ITimClient) {
	this.Client = client
}

type Conf struct {
	Heartbeat        int
	Name             string
	Pwd              string
	Domain           string
	Resource         string
	AckListener      func(ack *TimAckBean)
	MessageListener  func(mbean *TimMBean)
	PresenceListener func(pbean *TimPBean)

	MessageListListener  func(mbeans []*TimMBean)
	PresenceListListener func(pbeans []*TimPBean)
}

type Cli struct {
	Connect *Connect
	Sync    *sync.Mutex
	Flow    FLOW
	Addr    string
	conf    *Conf
	TLSAddr string
}

func NewCli(addr, tlsAddr string) (Timc *Cli) {
	Timc = &Cli{Sync: new(sync.Mutex), Flow: AUTH, Addr: addr, TLSAddr: tlsAddr}
	return
}

func (this *Cli) Sendmsg(toName string, msg *string) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("sendmsg,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.Sync.Lock()
	defer this.Sync.Unlock()
	mbean := new(TimMBean)
	mbean.Body = msg
	tid := new(Tid)
	domain, resource := this.conf.Domain, this.conf.Resource
	tid.Name = toName
	tid.Domain, tid.Resource = &domain, &resource
	mbean.ToTid = tid

	fid := new(Tid)
	fid.Name = this.conf.Name
	fid.Domain, fid.Resource = &domain, &resource
	mbean.FromTid = fid
	_type, maptype := "chat", int16(1)
	mbean.Type = &_type
	mbean.MsgType = &maptype
	return this.Connect.Client.TimMessage(mbean)
}

func (this *Cli) SendPresence(toName string, show *string) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SendPresence,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.Sync.Lock()
	defer this.Sync.Unlock()
	pbean := NewTimPBean()
	fromtid := NewTid()
	domain := this.conf.Domain
	fromtid.Domain = &domain
	fromtid.Name = this.conf.Name
	pbean.FromTid = fromtid

	totid := NewTid()
	totid.Domain = &domain
	totid.Name = toName
	pbean.ToTid = totid
	_type := "chat"
	pbean.Type = &_type
	pbean.Show = show
	return this.Connect.Client.TimPresence(pbean)
}

func (this *Cli) SendPBean(pbean *TimPBean) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SendPBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.Sync.Lock()
	defer this.Sync.Unlock()
	return this.Connect.Client.TimPresence(pbean)
}

func (this *Cli) SendMBean(mbean *TimMBean) {
	this.Connect.Client.TimMessage(mbean)
}

func (this *Cli) DisConnect() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("DisConnect,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if this != nil && this.Connect != nil {
		this.Connect.Client.Transport.Close()
	}
}

func (this *Cli) Ack(ab *TimAckBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Ack,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	this.Sync.Lock()
	defer this.Sync.Unlock()
	if this != nil && this.Connect != nil && this.Flow == AUTH {
		this.Connect.Client.TimAck(ab)
	}
}

func (this *Cli) Close() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Close,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if this != nil && this.Connect != nil {
		this.Flow = CLOSE
		this.Connect.Client.Transport.Close()
	}
}

func (this *Cli) Ping() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Ping,", err)
			logger.Error(string(debug.Stack()))
			this.Close()
			this = ReConn(this)
		}
	}()
	for {
		for i := 0; i < 20; i++ {
			time.Sleep(1 * time.Second)
		}
		if this == nil || this.Flow == CLOSE {
			break
		}
		if this != nil && this.Flow == AUTH {
			if this.Connect == nil {
				continue
			}
			err := func() (err error) {
				defer func() {
					if err := recover(); err != nil {
						logger.Error(string(debug.Stack()))
					}
				}()
				err = Ping(this)
				return
			}()
			if err != nil {
				logger.Error("ping error:", err.Error())
				break
			}
		}
	}
}

func Ping(this *Cli) (err error) {
	this.Sync.Lock()
	defer this.Sync.Unlock()
	err = this.Connect.Client.TimPing(fmt.Sprint(currentTimeMillis()))
	return
}

func getTransport(addr string) (transport *thrift.TSocket, err error) {
	transport, err = thrift.NewTSocket(addr)
	return
}

func getTSSLTransport(addr string) (transport *thrift.TSSLSocket, err error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport, err = thrift.NewTSSLSocket(addr, conf)
	return
}

func (this *Cli) Login() error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()

	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTCompactProtocolFactory()

	var useTransport thrift.TTransport
	if this.TLSAddr == "" {
		transport, err := getTransport(this.Addr)
		if err != nil {
			return err
		}
		if err := transport.Open(); err != nil {
			fmt.Fprintln(os.Stderr, "Error opening socket to ", this.Addr, " ", err)
			return err
		}
		useTransport = transportFactory.GetTransport(transport)
	} else {
		transport, err := getTSSLTransport(this.TLSAddr)
		if err != nil {
			return err
		}
		if err := transport.Open(); err != nil {
			fmt.Fprintln(os.Stderr, "Error opening socket to ", this.TLSAddr, " ", err)
			return err
		}
		useTransport = transportFactory.GetTransport(transport)
	}

	timclient := NewITimClientFactory(useTransport, protocolFactory)
	if this.Connect != nil {
		this.Connect.Close()
	}
	this.Connect = &Connect{FlowConnect: CONNECT_START}
	this.Connect.setITimClient(timclient)
	this.Connect.Super = this

	processorchan := make(chan int)
	go this.Connect.processor(processorchan)
	pro := <-processorchan
	if pro == 0 {
		return errors.New("processor error")
	}
	tid := new(Tid)
	domain, resource := this.conf.Domain, this.conf.Resource
	tid.Domain, tid.Resource = &domain, &resource
	tid.Name = this.conf.Name
	err := Login(this, tid)
	if err != nil {
		logger.Error("login err", err)
		return err
	}
	return nil
}

func Login(this *Cli, tid *Tid) (err error) {
	this.Sync.Lock()
	defer this.Sync.Unlock()
	param := NewTimParam()
	v, i := int16(Protocolversion), "1"
	param.Version = &v
	param.Interflow = &i
	this.Connect.Client.TimStream(param)
	err = this.Connect.Client.TimLogin(tid, this.conf.Pwd)
	return
}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func (this *Connect) processor(processorchan chan int) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
			if this.FlowConnect == CONNECT_START {
				processorchan <- 0
			}
		}
		if this != nil {
			this.FlowConnect = CONNECT_STOP
		}
	}()
	handler := new(TimImpl)
	handler.Client = this.Super
	processor := NewITimProcessor(handler)
	protocol := thrift.NewTCompactProtocol(this.Client.Transport)
	for {
		if this == nil || this.FlowConnect == CONNECT_STOP {
			break
		}
		if this.FlowConnect == CONNECT_START {
			this.FlowConnect = CONNECT_RUN
			processorchan <- 1
		}
		b, err := processor.Process(protocol, protocol)
		if err != nil && !b {
			break
		}
	}
}

func ReConn(cli *Cli) (client *Cli) {
	client, _ = NewConn(cli.Addr, cli.conf, cli.TLSAddr)
	return
}

func NewConn(addr string, conf *Conf, tlsAddr string) (cli *Cli, err error) {
	cli = NewCli(addr, tlsAddr)
	cli.conf = conf
	err = cli.Login()
	if err == nil {
		go cli.Ping()
	}
	return
}
