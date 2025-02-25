// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package timgo

import (
	"github.com/donnie4w/go-logger/logger"
	wss "github.com/donnie4w/gofer/websocket"
	"time"
)

type handle struct {
	tc        *TimClient
	pingCount int
	handler   *wss.Handler
}

func newHandle(tc *TimClient) (r *handle, err error) {
	r = &handle{tc: tc}
	if r.handler, err = wss.NewHandler(tc.cfg); err == nil {
		go r.ping()
	}
	return
}

func (h *handle) close() error {
	return h.handler.Close()
}

func (h *handle) pong() {
	h.pingCount = 0
}

func (h *handle) ping() {
	defer recoverable()
	ticker := time.NewTicker(15 * time.Second)
	for !h.tc.isClose {
		select {
		case <-ticker.C:
			h.pingCount++
			if h.tc.isClose {
				h.close()
				goto END
			}
			if err := h.send(h.tc.ts.ping()); err != nil || h.pingCount > 3 {
				logger.Error("ping over count>>", h.pingCount, err)
				h.close()
				goto END
			}
		}
	}
END:
}

func (h *handle) send(bs []byte) error {
	return h.handler.Send(bs)
}
