// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package cli

import (
	"fmt"
	"time"

	. "github.com/donnie4w/gofer/thrift"
	"github.com/donnie4w/simplelog/logging"
	. "github.com/donnie4w/timgo/stub"
)

type TimClient struct {
	addr            string
	pingCount       int
	handler         *cliHandle
	conf            *Config
	isLogout        bool
	Tx              *tx
	v               int
	messageHandler  func(*TimMessage)
	presenceHandler func(*TimPresence)
	streamHandler   func(*TimStream)
	nodesHandler    func(*TimNodes)
	ackHandler      func(*TimAck)

	pullmessageHandler   func(*TimMessageList)
	offlineMsgHandler    func(*TimMessageList)
	offlinemsgEndHandler func()

	bigStringHandler       func([]byte)
	bigBinaryHandler       func([]byte)
	bigBinaryStreamHandler func([]byte)
}

func NewTimClient(ip string, port int, tls bool) (tc *TimClient) {
	if addr := formatUrl(ip, port, tls); addr != "" {
		tc = &TimClient{addr: addr, Tx: &tx{&TimAuth{}}}
		tc.defaultInit()
	}
	return
}

func NewTimClientConfig(ip string, port int, tls bool, conf *Config) (tc *TimClient) {
	if addr := formatUrl(ip, port, tls); addr != "" {
		tc = &TimClient{addr: addr, Tx: &tx{&TimAuth{}}}
		conf.Url = addr
		tc.init(conf)
	}
	return
}

func (this *TimClient) login() (err error) {
	if this.handler != nil {
		err = this.handler.sendws(this.Tx.login())
	}
	return
}

func formatUrl(ip string, port int, tls bool) (url string) {
	if ip == "" || port > 65535 || port < 0 {
		return ""
	}
	if tls {
		url = fmt.Sprint("wss://", ip, ":", port)
	} else {
		url = fmt.Sprint("ws://", ip, ":", port)
	}
	return
}

func (this *TimClient) init(conf *Config) {
	this.conf = conf
	parse(this.conf)
	this.conf.OnError = func(_ *cliHandle, err error) {
		logging.Error("OnError:", err)
		this.close()
		<-time.After(time.Second << 2)
		if !this.isLogout {
			this.connect()
		}
	}
	this.conf.OnMessage = func(_ *cliHandle, msg []byte) {
		defer func() {
			if err := recover(); err != nil {
				logging.Error(err)
			}
		}()
		t := TIMTYPE(msg[0] & 0x7f)
		if msg[0]&0x80 == 0x80 {
			this.handler.sendws(append([]byte{byte(TIMACK)}, msg[1:5]...))
			msg = msg[5:]
		} else {
			msg = msg[1:]
		}
		this.pingCount = 0
		this.doMsg(t, msg)
	}
}

func (this *TimClient) defaultInit() {
	this.init(&Config{TimeOut: 10, Url: this.addr + "/tim", Origin: "https://github.com/donnie4w/tim"})
}

func (this *TimClient) Timeout(t time.Duration) {
	this.conf.TimeOut = t
}

func (this *TimClient) connect() (err error) {
	this.isLogout = false
	this.pingCount = 0
	this.v++
	if this.handler, err = NewCliHandle(this.conf); err == nil {
		this.login()
		go this.ping(this.v)
	} else {
		<-time.After(time.Second << 2)
		logging.Warn("reconn")
		this.connect()
	}
	<-time.After(time.Second)
	return
}

func (this *TimClient) CloseAndLogout() (err error) {
	this.isLogout = true
	this.v++
	return this.close()
}

func (this *TimClient) close() (err error) {
	if this.handler != nil {
		err = this.handler.Close()
	}
	return
}

func (this *TimClient) doMsg(t TIMTYPE, bs []byte) {
	switch t {
	case TIMPING:
		if this.pingCount > 0 {
			this.pingCount--
		}
	case TIMACK:
		if this.ackHandler != nil {
			if ta, err := TDecode(bs, &TimAck{}); ta != nil {
				this.ackHandler(ta)
			} else {
				logging.Error(err)
			}
		}
	case TIMMESSAGE:
		if this.messageHandler != nil {
			if tm, _ := TDecode(bs, &TimMessage{}); tm != nil {
				this.messageHandler(tm)
			}
		}
	case TIMPRESENCE:
		if this.presenceHandler != nil {
			if tp, _ := TDecode(bs, &TimPresence{}); tp != nil {
				this.presenceHandler(tp)
			}
		}
	case TIMNODES:
		if this.nodesHandler != nil {
			if tr, _ := TDecode(bs, &TimNodes{}); tr != nil {
				this.nodesHandler(tr)
			}
		}
	case TIMPULLMESSAGE:
		if this.pullmessageHandler != nil {
			if tm, _ := TDecode(bs, &TimMessageList{}); tm != nil {
				this.pullmessageHandler(tm)
			}
		}
	case TIMOFFLINEMSG:
		if this.offlineMsgHandler != nil {
			if tm, _ := TDecode(bs, &TimMessageList{}); tm != nil {
				this.offlineMsgHandler(tm)
			}
		}
	case TIMOFFLINEMSGEND:
		if this.offlinemsgEndHandler != nil {
			this.offlinemsgEndHandler()
		}
	case TIMSTREAM:
		if this.streamHandler != nil {
			if ts, _ := TDecode(bs, &TimStream{}); ts != nil {
				this.streamHandler(ts)
			}
		}
	case TIMBIGSTRING:
		if this.bigStringHandler != nil {
			this.bigStringHandler(bs)
		}
	case TIMBIGBINARY:
		if this.bigBinaryHandler != nil {
			this.bigBinaryHandler(bs)
		}
	case TIMBIGBINARYSTREAM:
		if this.bigBinaryStreamHandler != nil {
			this.bigBinaryStreamHandler(bs)
		}
	default:
		logging.Warn("undisposed >>>>>", t, " ,data length:", len(bs))
	}
}

func (this *TimClient) ping(v int) {
	defer _recover()
	ticker := time.NewTicker(15 * time.Second)
	for v == this.v {
		select {
		case <-ticker.C:
			this.pingCount++
			if v != this.v {
				goto END
			}
			if err := this.handler.sendws(this.Tx.ping()); err != nil || this.pingCount > 3 {
				logging.Error("ping over count>>", this.pingCount, err)
				this.close()
				goto END
			}
		}
	}
END:
}

func (this *TimClient) MessageHandler(handler func(*TimMessage)) {
	this.messageHandler = handler
}

func (this *TimClient) PresenceHandler(handler func(*TimPresence)) {
	this.presenceHandler = handler
}

func (this *TimClient) AckHandler(handler func(*TimAck)) {
	this.ackHandler = handler
}

func (this *TimClient) NodesHandler(handler func(*TimNodes)) {
	this.nodesHandler = handler
}

func (this *TimClient) PullmessageHandler(handler func(*TimMessageList)) {
	this.pullmessageHandler = handler
}

func (this *TimClient) OfflineMsgHandler(handler func(*TimMessageList)) {
	this.offlineMsgHandler = handler
}

func (this *TimClient) OfflinemsgEndHandler(handler func()) {
	this.offlinemsgEndHandler = handler
}

func (this *TimClient) StreamHandler(handler func(*TimStream)) {
	this.streamHandler = handler
}
