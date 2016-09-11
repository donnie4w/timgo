package client

import (
	"runtime/debug"

	"github.com/donnie4w/go-logger/logger"
	. "timgo.protocol"
)

type TimImpl struct {
	Ip     string
	Port   int
	Pub    string //发布id
	Client *Cli
}

// Parameters:
//  - Param
func (this *TimImpl) TimStream(param *TimParam) (err error) {
	return
}
func (this *TimImpl) TimStarttls() (err error) {
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimLogin(tid *Tid, pwd string) (err error) {
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
func (this *TimImpl) TimAck(ab *TimAckBean) (err error) {
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
func (this *TimImpl) TimPresence(pbean *TimPBean) (err error) {
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
func (this *TimImpl) TimMessage(mbean *TimMBean) (err error) {
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
func (this *TimImpl) TimPing(threadId string) (err error) {
	ab := new(TimAckBean)
	ab.ID = &threadId
	acktype, ackstatus := "ping", "200"
	ab.AckType, ab.AckStatus = &acktype, &ackstatus
	this.Client.Ack(ab)
	return
}

// Parameters:
//  - E
func (this *TimImpl) TimError(e *TimError) (err error) {
	return
}
func (this *TimImpl) TimLogout() (err error) {
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimRegist(tid *Tid, pwd string) (err error) {
	return
}

// Parameters:
//  - Tid
//  - Pwd
func (this *TimImpl) TimRemoteUserAuth(tid *Tid, pwd string, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	return
}

// Parameters:
//  - Tid
func (this *TimImpl) TimRemoteUserGet(tid *Tid, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	return
}

// Parameters:
//  - Tid
//  - Ub
func (this *TimImpl) TimRemoteUserEdit(tid *Tid, ub *TimUserBean, auth *TimAuth) (r *TimRemoteUserBean, err error) {
	return
}

// Parameters:
//  - Pbean
func (this *TimImpl) TimResponsePresence(pbean *TimPBean, auth *TimAuth) (r *TimResponseBean, err error) {
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimResponseMessage(mbean *TimMBean, auth *TimAuth) (r *TimResponseBean, err error) {
	logger.Debug("ResponseMessage", mbean)
	return
}

func (this *TimImpl) TimMessageIq(timMsgIq *TimMessageIq, iqType string) (err error) {
	logger.Debug("TimMessageIq:", timMsgIq, " ", iqType)
	return
}

// Parameters:
//  - Mbean
func (this *TimImpl) TimMessageResult_(mbean *TimMBean) (err error) {
	logger.Debug("TimMessageResult_:", mbean)
	return
}

func (this *TimImpl) TimRoser(roster *TimRoster) (err error) {
	logger.Debug("TimRoser:", roster)
	return
}

func (this *TimImpl) TimResponseMessageIq(timMsgIq *TimMessageIq, iqType string, auth *TimAuth) (r *TimMBeanList, err error) {
	logger.Debug("TimResponseMessageIq:", timMsgIq, iqType, auth)
	panic("error TimResponseMessageIq")
	return
}

func (this *TimImpl) TimMessageList(mbeanList *TimMBeanList) (err error) {
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
func (this *TimImpl) TimPresenceList(pbeanList *TimPBeanList) (err error) {
	ab := NewTimAckBean()
	id := pbeanList.ThreadId
	ab.ID = &id
	ackstatus, acktype := "200", "presence"
	ab.AckStatus, ab.AckType = &ackstatus, &acktype
	this.Client.Ack(ab)
	this.Client.conf.PresenceListListener(pbeanList.GetTimPBeanList())
	return
}

func (this *TimImpl) TimResponsePresenceList(pbeanList *TimPBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	panic("error TimResponsePresenceList")
	return
}

// Parameters:
//  - MbeanList
//  - Auth
func (this *TimImpl) TimResponseMessageList(mbeanList *TimMBeanList, auth *TimAuth) (r *TimResponseBean, err error) {
	panic("error TimResponseMessageList")
	return
}

func (this *TimImpl) TimProperty(tpb *TimPropertyBean) (err error) {
	logger.Debug("TimProperty:", tpb)
	return
}
