package client

import (
	"context"
	"runtime/debug"

	. "timgo/protocol"

	"github.com/donnie4w/go-logger/logger"
)

type TimImpl struct {
	Ip     string
	Port   int
	Pub    string //发布id
	Client *Cli
}

// Parameters:
//  - Param
func (this *TimImpl) TimStream(ctx context.Context, param *TimParam) (err error) {
	return
}
func (this *TimImpl) TimStarttls(ctx context.Context) (err error) {
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimLogin(ctx context.Context, tid *Tid, pwd string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Login error", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("Login:", tid, pwd)
	return
}

// Parameters:
//  - Ab
func (this *TimImpl) TimAck(ctx context.Context, ab *TimAckBean) (err error) {
	if ab != nil {
		switch *ab.AckType {
		case "login":
			switch *ab.AckStatus {
			case "200":
				this.Client.Flow = AUTH
			case "400":
				this.Client.Flow = NOAUTH
			}
		default:
			this.Client.conf.AckListener(ab)
		}
	}
	return
}

// Parameters:
//  - Pbean
func (this *TimImpl) TimPresence(ctx context.Context, pbean *TimPBean) (err error) {
	id := pbean.ThreadId
	ab := NewTimAckBean()
	ab.ID = &id
	ackstatus, acktype := "200", "presence"
	ab.AckStatus, ab.AckType = &ackstatus, &acktype
	this.Client.Ack(ab)
	if pbean.GetStatus() == "probe" {
		p := NewTimPBean()
		p.FromTid = NewTid()
		p.FromTid.Domain = pbean.FromTid.Domain
		p.FromTid.Name = this.Client.conf.Name
		p.ToTid = NewTid()
		p.ToTid.Domain = pbean.FromTid.Domain
		p.ToTid.Name = pbean.GetFromTid().GetName()
		status, show := "available", "online"
		p.Status, p.Show = &status, &show
		this.Client.SendPBean(p)
	}
	this.Client.conf.PresenceListener(pbean)
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimMessage(ctx context.Context, mbean *TimMBean) (err error) {
	ab := NewTimAckBean()
	id := mbean.ThreadId
	ab.ID = &id
	ackstatus, acktype := "200", "message"
	ab.AckStatus, ab.AckType = &ackstatus, &acktype
	this.Client.Ack(ab)
	this.Client.conf.MessageListener(mbean)
	return
}

// Parameters:
//  - ThreadId
func (this *TimImpl) TimPing(ctx context.Context, threadId string) (err error) {
	ab := new(TimAckBean)
	ab.ID = &threadId
	acktype, ackstatus := "ping", "200"
	ab.AckType, ab.AckStatus = &acktype, &ackstatus
	this.Client.Ack(ab)
	return
}

// Parameters:
//  - E
func (this *TimImpl) TimError(ctx context.Context, e *TimError) (err error) {
	return
}
func (this *TimImpl) TimLogout(ctx context.Context) (err error) {
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimRegist(ctx context.Context, tid *Tid, pwd string) (err error) {
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimRemoteUserAuth(ctx context.Context, tid *Tid, pwd string, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	return
}

// Parameters:
//  - Tid
func (this *TimImpl) TimRemoteUserGet(ctx context.Context, tid *Tid, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	return
}

// Parameters:
//  - Tid
//  - Ub
func (this *TimImpl) TimRemoteUserEdit(ctx context.Context, tid *Tid, ub *TimUserBean, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	return
}

// Parameters:
//  - Pbean
func (this *TimImpl) TimResponsePresence(ctx context.Context, pbean *TimPBean, auth *TimAuth) (r *TimResponseBean, err error) {
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimResponseMessage(ctx context.Context, mbean *TimMBean, auth *TimAuth) (r *TimResponseBean, err error) {
	logger.Debug("ResponseMessage", mbean)
	return
}

func (this *TimImpl) TimMessageIq(ctx context.Context, timMsgIq *TimMessageIq, iqType string) (err error) {
	logger.Debug("TimMessageIq:", timMsgIq, " ", iqType)
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimMessageResult_(ctx context.Context, mbean *TimMBean) (err error) {
	logger.Debug("TimMessageResult_:", mbean)
	return
}

func (this *TimImpl) TimRoser(ctx context.Context, roster *TimRoster) (err error) {
	logger.Debug("TimRoser:", roster)
	return
}

func (this *TimImpl) TimResponseMessageIq(ctx context.Context, timMsgIq *TimMessageIq, iqType string, auth *TimAuth) (r *TimMBeanList, err error) {
	logger.Debug("TimResponseMessageIq:", timMsgIq, iqType, auth)
	panic("error TimResponseMessageIq")
}

func (this *TimImpl) TimMessageList(ctx context.Context, mbeanList *TimMBeanList) (err error) {
	ab := NewTimAckBean()
	id := mbeanList.ThreadId
	ab.ID = &id
	ackstatus, acktype := "200", "message"
	ab.AckStatus, ab.AckType = &ackstatus, &acktype
	this.Client.Ack(ab)
	this.Client.conf.MessageListListener(mbeanList.GetTimMBeanList())
	return
}

// Parameters:
//  - PbeanList
func (this *TimImpl) TimPresenceList(ctx context.Context, pbeanList *TimPBeanList) (err error) {
	ab := NewTimAckBean()
	id := pbeanList.ThreadId
	ab.ID = &id
	ackstatus, acktype := "200", "presence"
	ab.AckStatus, ab.AckType = &ackstatus, &acktype
	this.Client.Ack(ab)
	this.Client.conf.PresenceListListener(pbeanList.GetTimPBeanList())
	return
}

func (this *TimImpl) TimResponsePresenceList(pctx context.Context, beanList *TimPBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	panic("error TimResponsePresenceList")
}

// Parameters:
//  - MbeanList
//  - Auth
func (this *TimImpl) TimResponseMessageList(ctx context.Context, mbeanList *TimMBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	panic("error TimResponseMessageList")
}

func (this *TimImpl) TimProperty(ctx context.Context, tpb *TimPropertyBean) (err error) {
	logger.Debug("TimProperty:", tpb)
	return
}
