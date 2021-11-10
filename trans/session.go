package trans

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"time"
)

type SessionStatus int

const (
	Working SessionStatus = 0
	// WaitingTarget SessionStatus = 1
	// WaitingClient SessionStatus = 2
)

type SessionObject struct {
	id     string
	target *DeviceObject
	client *DeviceObject
	status SessionStatus
}

func SessionObjectNew(target *DeviceObject, client *DeviceObject) (string, *SessionObject) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", nil
	}
	target.SetStatus(Busy)
	client.SetStatus(Busy)
	idStr := id.String()
	return idStr, &SessionObject{
		id:     idStr,
		target: target,
		client: client,
		status: Working,
	}
}

// 成员状态广播

// ForwardToClient 转发 target->client
/*
func (m *SessionObject)ForwardToClient(messageType int, data []byte) error {
	return m.client.client.WriteMessage(messageType, data)
}
*/
// ForwardToTarget 转发 client->target
/*
func (m *SessionObject)ForwardToTarget(messageType int, data []byte) error {
	return m.target.client.WriteMessage(messageType, data)
}
*/

func (m *SessionObject) RestoreTarget(o *websocket.Conn) {
	m.target.client = o
	m.target.status = Busy
	m.toClient(TargetReconnect, nil)
	for m.target.work {
		time.Sleep(time.Duration(2) * time.Second)
	}
	go func() {
		err := ForwardToClient(m)
		if err != nil {
			panic(err)
		}
	}()
}

func (m *SessionObject) RestoreClient(o *websocket.Conn) {
	m.client.client = o
	m.client.status = Busy
	for m.client.work {
		time.Sleep(time.Duration(2) * time.Second)
	}
	go func() {
		err := ForwardToTarget(m)
		if err != nil {
			panic(err)
		}
	}()
}

func ForwardToTarget(m *SessionObject) error {
	var e error = nil
	if m.client.work == false {
		m.client.work = true
		for m.client.status == Busy {
			// m.client.SetReadDeadline( time.Now().Add(time.Duration(5) * time.Second) )
			t, c, e := m.client.client.ReadMessage()
			if t == -1 || e != nil {
				break
			}
			e = m.client.GetSession().target.client.WriteMessage(t, c)
			if e != nil {
				break
			}
		}
		m.client.work = false
	}
	return e
}

func ForwardToClient(m *SessionObject) error {
	var e error = nil
	if m.target.work == false {
		m.target.work = true
		for m.target.status == Busy {
			t, c, e := m.target.client.ReadMessage()
			if t == -1 || e != nil {
				break
			}
			e = m.target.GetSession().client.client.WriteMessage(t, c)
			if e != nil {
				break
			}
		}
		m.target.work = false
	}
	return e
}

func (m *SessionObject) Working() {
	// m.client.SetStatus(Busy)
	go func() {
		err := ForwardToTarget(m)
		if err != nil {
			panic(err)
		}
	}()
	// m.target.SetStatus(Busy)
	go func() {
		err := ForwardToClient(m)
		if err != nil {
			panic(err)
		}
	}()
}

func (m *SessionObject) toTarget(cmd resultCmd, content interface{}) {
	r := RealMessageNew(Success, cmd, content)
	err := m.target.client.WriteJSON(r)
	if err != nil {
		panic(err)
	}
}
func (m *SessionObject) toClient(cmd resultCmd, content interface{}) {
	r := RealMessageNew(Success, cmd, content)
	err := m.client.client.WriteJSON(r)
	if err != nil {
		panic(err)
	}
}

func (m *SessionObject) BroadCast(cmd resultCmd, content interface{}) {
	r := RealMessageNew(Success, cmd, content)
	err := m.client.client.WriteJSON(r)
	if err != nil {
		panic(err)
	}
	err = m.target.client.WriteJSON(r)
	if err != nil {
		panic(err)
	}
}

func (m *SessionObject) CanDestroy() bool {
	return m.client.status == Lost && m.target.status == Lost
}

func (m *SessionObject) GetId() string {
	return m.id
}

func (m *SessionObject) GetTarget() *DeviceObject {
	return m.target
}
