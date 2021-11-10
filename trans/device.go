package trans

import (
	"github.com/gorilla/websocket"
)

type DeviceStatus int
const (
	Online	DeviceStatus = 0
	Lost	DeviceStatus = 1
	Busy	DeviceStatus = 2 // 会话中
	Ban		DeviceStatus = 4
)

type DeviceObject struct{
	name string
	status DeviceStatus					// 设备状态
	client *websocket.Conn			// ws客户端
	session *SessionObject
	work bool
}

func DeviceObjectNew(name string, client *websocket.Conn) *DeviceObject{
	return &DeviceObject{
		name: name,
		status: Online,
		client: client,
		session: nil,
		work: false,
	}
}

func (m *DeviceObject) SetStatus(status DeviceStatus) {
	m.status = status
}

func (m *DeviceObject) GetStatus() DeviceStatus {
	return m.status
}

func (m *DeviceObject) JoinSession(session *SessionObject) {
	m.session = session
}

func (m *DeviceObject) LeaveSession() {
	if m.session != nil {
		m.status = Online
		m.session = nil
	}
	// 等待 work 为 false 才说明退出成功
	/*
	for m.work {
		time.Sleep(time.Duration(2) * time.Second)
	}
	*/
}

func (m *DeviceObject) GetSession() *SessionObject{
	return m.session
}