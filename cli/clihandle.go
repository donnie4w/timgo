// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/timgo

package cli

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/donnie4w/simplelog/logging"
	"golang.org/x/net/websocket"
)

type Config struct {
	Url       string
	Origin    string
	HttpUrl   string
	CertFiles []string
	CertBytes [][]byte
	TimeOut   time.Duration
	Tx        *tx
	OnOpen    func(c *cliHandle)
	OnError   func(c *cliHandle, err error)
	OnClose   func(c *cliHandle)
	OnMessage func(c *cliHandle, msg []byte)
}

type cliHandle struct {
	conf     *Config
	conn     *websocket.Conn
	mux      *sync.Mutex
	_isError bool
	_isAuth  bool
}

func NewCliHandle(conf *Config) (cli *cliHandle, err error) {
	var conn *websocket.Conn
	config := &websocket.Config{Dialer: &net.Dialer{Timeout: conf.TimeOut * time.Second}, Version: websocket.ProtocolVersionHybi13}
	if strings.HasPrefix(conf.Url, "wss:") {
		if conf.CertFiles != nil {
			if rootcas, err := loadCACertificatesFromFiles(conf.CertFiles); err == nil {
				config.TlsConfig = &tls.Config{
					RootCAs:            rootcas,
					InsecureSkipVerify: false,
				}
			} else {
				return nil, err
			}
		} else if conf.CertBytes != nil {
			if rootcas, err := loadCACertificatesFromBytes(conf.CertBytes); err == nil {
				config.TlsConfig = &tls.Config{
					RootCAs:            rootcas,
					InsecureSkipVerify: false,
				}
			} else {
				return nil, err
			}
		} else {
			config.TlsConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
	if config.Location, err = url.ParseRequestURI(conf.Url); err == nil {
		if config.Origin, err = url.ParseRequestURI(conf.Origin); err == nil {
			conn, err = websocket.DialConfig(config)
		}
	}
	if err == nil && conn != nil {
		cli = &cliHandle{conf, conn, &sync.Mutex{}, false, false}
		if conf.OnOpen != nil {
			conf.OnOpen(cli)
		}
		go cli._read()
	} else {
		logging.Error(err)
	}
	return
}

func (this *cliHandle) _sendws(bs []byte) (err error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	return websocket.Message.Send(this.conn, bs)
}

func (this *cliHandle) sendws(bs []byte) (err error) {
	return this._sendws(bs)
}

func (this *cliHandle) Close() (err error) {
	if this.conn != nil {
		this._isError = true
		err = this.conn.Close()
	}
	return
}

func (this *cliHandle) _read() {
	var err error
	for !this._isError {
		var byt []byte
		if err = websocket.Message.Receive(this.conn, &byt); err != nil {
			this._isError = true
			break
		}
		if byt != nil && this.conf.OnMessage != nil {
			go this.conf.OnMessage(this, byt)
		}
	}
	if this.conf.OnError != nil {
		go this.conf.OnError(this, err)
	}
	this.Close()
	if this.conf.OnClose != nil {
		this.conf.OnClose(this)
	}
}

func _recover() {
	if err := recover(); err != nil {
		logging.Error(err)
	}
}

/***********************************************/
func parse(conf *Config) {
	ss := strings.Split(conf.Url, "//")
	s := strings.Split(ss[1], "/")
	url := "http"
	if strings.HasPrefix(ss[0], "wss:") {
		url = "https"
	}
	url = url + "://" + s[0] + "/tim2"
	conf.HttpUrl = url
}

func sendsync(conf *Config, bs []byte) (_r []byte, err error) {
	_r, err = httpPost(bs, conf, true)
	return
}

func httpPost(bs []byte, conf *Config, close bool) (_r []byte, err error) {
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	if strings.HasPrefix(conf.HttpUrl, "https:") {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := http.Client{Transport: tr}
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, conf.HttpUrl, bytes.NewReader(bs)); err == nil {
		if close {
			req.Close = true
		}
		req.Header.Set("Origin", conf.Origin)
		var resp *http.Response
		if resp, err = client.Do(req); err == nil {
			if close {
				defer resp.Body.Close()
			}
			var body []byte
			if body, err = io.ReadAll(resp.Body); err == nil {
				_r = body
			}
		}
	}
	return
}

func loadCACertificatesFromFiles(certFiles []string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, certFile := range certFiles {
		pemData, err := os.ReadFile(certFile)
		if err != nil {
			return nil, err
		}

		if ok := pool.AppendCertsFromPEM(pemData); !ok {
			return nil, errors.New("failed to append certificates from PEM data")
		}
	}
	return pool, nil
}

func loadCACertificatesFromBytes(certBytes [][]byte) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, bs := range certBytes {
		if ok := pool.AppendCertsFromPEM(bs); !ok {
			return nil, errors.New("failed to append certificates from PEM data")
		}
	}
	return pool, nil
}
