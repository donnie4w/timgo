// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package timgo

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/donnie4w/gofer/httpclient"
	"github.com/donnie4w/gofer/thrift"
	wss "github.com/donnie4w/gofer/websocket"
	"github.com/donnie4w/timgo/stub"
	"runtime/debug"
	"strings"
	"time"
)

type TimClient struct {
	addr                   string
	handler                *handle
	cfg                    *wss.Config
	pongTime               int64
	isClose                bool
	ts                     *tx
	syncUrl                string
	messageHandler         func(*stub.TimMessage)
	presenceHandler        func(*stub.TimPresence)
	streamHandler          func(*stub.TimStream)
	nodesHandler           func(*stub.TimNodes)
	ackHandler             func(*stub.TimAck)
	pullmessageHandler     func(*stub.TimMessageList)
	offlineMsgHandler      func(*stub.TimMessageList)
	offlinemsgEndHandler   func()
	bigStringHandler       func([]byte)
	bigBinaryHandler       func([]byte)
	bigBinaryStreamHandler func([]byte)
}

func NewTimClient(tls bool, ip string, port int) (tc *TimClient) {
	if addr := formatUrl(ip, port, tls); addr != "" {
		tc = &TimClient{addr: addr, ts: &tx{&stub.TimAuth{}}, pongTime: time.Now().UnixNano()}
		tc.defaultInit()
		go tc.reconnect()
	}
	return
}

func NewTimClientWithConfig(ip string, port int, tls bool, conf *wss.Config) (tc *TimClient) {
	if addr := formatUrl(ip, port, tls); addr != "" {
		tc = &TimClient{addr: addr, ts: &tx{&stub.TimAuth{}}, pongTime: time.Now().UnixNano()}
		conf.Url = addr
		tc.init(conf)
		go tc.reconnect()
	}
	return
}

func (tc *TimClient) login() (err error) {
	if tc.handler != nil {
		err = tc.handler.send(tc.ts.login())
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
	tc.cfg.OnError = func(handle *wss.Handler, err error) {
		logger.Error("OnError:", err)
		handle.Close()
	}
	tc.cfg.OnMessage = func(handle *wss.Handler, msg []byte) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err)
			}
		}()
		tc.pongTime = time.Now().UnixNano()
		tc.handler.pong()
		t := TIMTYPE(msg[0] & 0x7f)
		if msg[0]&0x80 == 0x80 {
			handle.Send(append([]byte{byte(TIMACK)}, msg[1:5]...))
			msg = msg[5:]
		} else {
			msg = msg[1:]
		}
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
	if tc.handler != nil {
		tc.handler.close()
	}
	if tc.handler, err = newHandle(tc); err == nil {
		tc.login()
	}
	return
}

func (tc *TimClient) close() (err error) {
	tc.isClose = true
	return tc.closeconnect()
}

func (tc *TimClient) closeconnect() (err error) {
	if tc.handler != nil {
		err = tc.handler.close()
	}
	return
}

func (tc *TimClient) doMsg(t TIMTYPE, bs []byte) {
	switch t {
	case TIMPING:
	case TIMACK:
		if tc.ackHandler != nil {
			if ta, err := thrift.TDecode(bs, &stub.TimAck{}); ta != nil {
				tc.ackHandler(ta)
			} else {
				logger.Error(err)
			}
		}
	case TIMMESSAGE:
		if tc.messageHandler != nil {
			if tm, _ := thrift.TDecode(bs, &stub.TimMessage{}); tm != nil {
				tc.messageHandler(tm)
			}
		}
	case TIMPRESENCE:
		if tc.presenceHandler != nil {
			if tp, _ := thrift.TDecode(bs, &stub.TimPresence{}); tp != nil {
				tc.presenceHandler(tp)
			}
		}
	case TIMNODES:
		if tc.nodesHandler != nil {
			if tr, _ := thrift.TDecode(bs, &stub.TimNodes{}); tr != nil {
				tc.nodesHandler(tr)
			}
		}
	case TIMPULLMESSAGE:
		if tc.pullmessageHandler != nil {
			if tm, _ := thrift.TDecode(bs, &stub.TimMessageList{}); tm != nil {
				tc.pullmessageHandler(tm)
			}
		}
	case TIMOFFLINEMSG:
		if tc.offlineMsgHandler != nil {
			if tm, _ := thrift.TDecode(bs, &stub.TimMessageList{}); tm != nil {
				tc.offlineMsgHandler(tm)
			}
		}
	case TIMOFFLINEMSGEND:
		if tc.offlinemsgEndHandler != nil {
			tc.offlinemsgEndHandler()
		}
	case TIMSTREAM:
		if tc.streamHandler != nil {
			if ts, _ := thrift.TDecode(bs, &stub.TimStream{}); ts != nil {
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

//func (tc *TimClient) ping() {
//	defer recoverable()
//	ticker := time.NewTicker(15 * time.Second)
//	for !tc.isClose {
//		select {
//		case <-ticker.C:
//			tc.pingCount++
//			if tc.isClose {
//				goto END
//			}
//			if err := tc.handler.Send(tc.ts.ping()); err != nil || tc.pingCount > 3 {
//				logger.Error("ping over count>>", tc.pingCount, err)
//				tc.closeconnect()
//				goto END
//			}
//		}
//	}
//END:
//}

func (tc *TimClient) MessageHandler(handler func(*stub.TimMessage)) {
	tc.messageHandler = handler
}

func (tc *TimClient) PresenceHandler(handler func(*stub.TimPresence)) {
	tc.presenceHandler = handler
}

func (tc *TimClient) AckHandler(handler func(*stub.TimAck)) {
	tc.ackHandler = handler
}

func (tc *TimClient) NodesHandler(handler func(*stub.TimNodes)) {
	tc.nodesHandler = handler
}

func (tc *TimClient) PullmessageHandler(handler func(*stub.TimMessageList)) {
	tc.pullmessageHandler = handler
}

func (tc *TimClient) OfflineMsgHandler(handler func(*stub.TimMessageList)) {
	tc.offlineMsgHandler = handler
}

func (tc *TimClient) OfflineMsgEndHandler(handler func()) {
	tc.offlinemsgEndHandler = handler
}

func (tc *TimClient) StreamHandler(handler func(*stub.TimStream)) {
	tc.streamHandler = handler
}

func recoverable() {
	if err := recover(); err != nil {
		logger.Error(string(debug.Stack()))
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

func (tc *TimClient) reconnect() {
	defer recoverable()
	ticker := time.NewTicker(15 * time.Second)
	for !tc.isClose {
		select {
		case <-ticker.C:
			if time.Now().UnixNano()-tc.pongTime > int64(time.Minute) {
				tc.connect()
			}
		}
	}
}
