// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package timgo

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/donnie4w/gofer/httpclient"
	. "github.com/donnie4w/gofer/thrift"
	wss "github.com/donnie4w/gofer/websocket"
	. "github.com/donnie4w/timgo/stub"
	"strings"
	"time"
)

type TimClient struct {
	addr                   string
	pingCount              int
	handler                *wss.Handler
	cfg                    *wss.Config
	isClose                bool
	ts                     *tx
	syncUrl                string
	messageHandler         func(*TimMessage)
	presenceHandler        func(*TimPresence)
	streamHandler          func(*TimStream)
	nodesHandler           func(*TimNodes)
	ackHandler             func(*TimAck)
	pullmessageHandler     func(*TimMessageList)
	offlineMsgHandler      func(*TimMessageList)
	offlinemsgEndHandler   func()
	bigStringHandler       func([]byte)
	bigBinaryHandler       func([]byte)
	bigBinaryStreamHandler func([]byte)
	AfterLoginEvent        func()
}

func NewTimClient(tls bool, ip string, port int) (tc *TimClient) {
	if addr := formatUrl(ip, port, tls); addr != "" {
		tc = &TimClient{addr: addr, ts: &tx{&TimAuth{}}}
		tc.defaultInit()
	}
	return
}

func NewTimClientWithConfig(ip string, port int, tls bool, conf *wss.Config) (tc *TimClient) {
	if addr := formatUrl(ip, port, tls); addr != "" {
		tc = &TimClient{addr: addr, ts: &tx{&TimAuth{}}}
		conf.Url = addr
		tc.init(conf)
	}
	return
}

func (tc *TimClient) login() (err error) {
	if tc.handler != nil {
		err = tc.handler.Send(tc.ts.login())
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

func (tc *TimClient) init(config *wss.Config) {
	tc.cfg = config
	tc.syncUrl = parse(tc.cfg)
	tc.cfg.OnError = func(_ *wss.Handler, err error) {
		logger.Error("OnError:", err)
		tc.closeconnect()
		<-time.After(time.Second << 2)
		if !tc.isClose {
			tc.connect()
		}
	}
	tc.cfg.OnMessage = func(_ *wss.Handler, msg []byte) {
		defer func() {
			if e := recover(); e != nil {
				logger.Error(e)
			}
		}()
		t := TIMTYPE(msg[0] & 0x7f)
		if msg[0]&0x80 == 0x80 {
			tc.handler.Send(append([]byte{byte(TIMACK)}, msg[1:5]...))
			msg = msg[5:]
		} else {
			msg = msg[1:]
		}
		tc.pingCount = 0
		tc.doMsg(t, msg)
	}
}

func (tc *TimClient) defaultInit() {
	tc.init(&wss.Config{TimeOut: 15 * time.Second, Url: tc.addr + "/tim", Origin: "https://github.com/donnie4w/tim"})
}

func (tc *TimClient) TimeOut(t time.Duration) {
	tc.cfg.TimeOut = t
}

func (tc *TimClient) connect() (err error) {
	tc.isClose = false
	tc.pingCount = 0
	if tc.handler, err = wss.NewHandler(tc.cfg); err == nil {
		if err = tc.login(); err == nil {
			go tc.ping()
		}
	} else {
		fmt.Println(err)
		<-time.After(4 * time.Second)
		logger.Warn("reconnect")
		tc.connect()
	}
	<-time.After(time.Second)
	return
}

func (tc *TimClient) close() (err error) {
	tc.isClose = true
	return tc.closeconnect()
}

func (tc *TimClient) closeconnect() (err error) {
	if tc.handler != nil {
		err = tc.handler.Close()
	}
	return
}

func (tc *TimClient) doMsg(t TIMTYPE, bs []byte) {
	switch t {
	case TIMPING:
		if tc.pingCount > 0 {
			tc.pingCount--
		}
	case TIMACK:
		if tc.ackHandler != nil {
			if ta, err := TDecode(bs, &TimAck{}); ta != nil {
				tc.ackHandler(ta)
			} else {
				logger.Error(err)
			}
		}
	case TIMMESSAGE:
		if tc.messageHandler != nil {
			if tm, _ := TDecode(bs, &TimMessage{}); tm != nil {
				tc.messageHandler(tm)
			}
		}
	case TIMPRESENCE:
		if tc.presenceHandler != nil {
			if tp, _ := TDecode(bs, &TimPresence{}); tp != nil {
				tc.presenceHandler(tp)
			}
		}
	case TIMNODES:
		if tc.nodesHandler != nil {
			if tr, _ := TDecode(bs, &TimNodes{}); tr != nil {
				tc.nodesHandler(tr)
			}
		}
	case TIMPULLMESSAGE:
		if tc.pullmessageHandler != nil {
			if tm, _ := TDecode(bs, &TimMessageList{}); tm != nil {
				tc.pullmessageHandler(tm)
			}
		}
	case TIMOFFLINEMSG:
		if tc.offlineMsgHandler != nil {
			if tm, _ := TDecode(bs, &TimMessageList{}); tm != nil {
				tc.offlineMsgHandler(tm)
			}
		}
	case TIMOFFLINEMSGEND:
		if tc.offlinemsgEndHandler != nil {
			tc.offlinemsgEndHandler()
		}
	case TIMSTREAM:
		if tc.streamHandler != nil {
			if ts, _ := TDecode(bs, &TimStream{}); ts != nil {
				tc.streamHandler(ts)
			}
		}
	case TIMBIGSTRING:
		if tc.bigStringHandler != nil {
			tc.bigStringHandler(bs)
		}
	case TIMBIGBINARY:
		if tc.bigBinaryHandler != nil {
			tc.bigBinaryHandler(bs)
		}
	case TIMBIGBINARYSTREAM:
		if tc.bigBinaryStreamHandler != nil {
			tc.bigBinaryStreamHandler(bs)
		}
	default:
		logger.Warn("undisposed >>>>>", t, " ,data length:", len(bs))
	}
}

func (tc *TimClient) ping() {
	defer recoverable()
	ticker := time.NewTicker(15 * time.Second)
	for !tc.isClose {
		select {
		case <-ticker.C:
			tc.pingCount++
			if tc.isClose {
				goto END
			}
			if err := tc.handler.Send(tc.ts.ping()); err != nil || tc.pingCount > 3 {
				logger.Error("ping over count>>", tc.pingCount, err)
				tc.closeconnect()
				goto END
			}
		}
	}
END:
}

func (tc *TimClient) MessageHandler(handler func(*TimMessage)) {
	tc.messageHandler = handler
}

func (tc *TimClient) PresenceHandler(handler func(*TimPresence)) {
	tc.presenceHandler = handler
}

func (tc *TimClient) AckHandler(handler func(*TimAck)) {
	tc.ackHandler = handler
}

func (tc *TimClient) NodesHandler(handler func(*TimNodes)) {
	tc.nodesHandler = handler
}

func (tc *TimClient) PullmessageHandler(handler func(*TimMessageList)) {
	tc.pullmessageHandler = handler
}

func (tc *TimClient) OfflineMsgHandler(handler func(*TimMessageList)) {
	tc.offlineMsgHandler = handler
}

func (tc *TimClient) OfflineMsgEndHandler(handler func()) {
	tc.offlinemsgEndHandler = handler
}

func (tc *TimClient) StreamHandler(handler func(*TimStream)) {
	tc.streamHandler = handler
}

func recoverable() {
	if err := recover(); err != nil {
		logger.Error(err)
	}
}

func parse(cfg *wss.Config) string {
	ss := strings.Split(cfg.Url, "//")
	s := strings.Split(ss[1], "/")
	url := "http"
	if strings.HasPrefix(ss[0], "wss:") {
		url = "https"
	}
	return url + "://" + s[0] + "/tim2"
}

func sendsync(url, origin string, bs []byte) ([]byte, error) {
	m := map[string]string{}
	if origin != "" {
		m["Origin"] = origin
	}
	return httpclient.Post3(bs, true, url, m, nil)
}
