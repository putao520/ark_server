package trans

import (
	"errors"
)

type SessionManager struct {
	sessionQueue map[string]*SessionObject
}

var sessionManagerInstance *SessionManager
func SessionManagerInstance() *SessionManager{
	if sessionManagerInstance == nil {
		sessionManagerInstance = SessionManagerNew()
	}
	return sessionManagerInstance
}

func SessionManagerNew() *SessionManager{
	return &SessionManager{
		sessionQueue: make(map[string]*SessionObject),
	}
}

func (m *SessionManager) Build(target *DeviceObject, client *DeviceObject) *SessionObject{
	id, obj := SessionObjectNew(target, client)
	if len(id) == 0 {
		panic("uuid create failed")
	}
	m.sessionQueue[id] = obj
	target.JoinSession(obj)
	client.JoinSession(obj)
	obj.Working()
	return obj
}

func (m *SessionManager) Get(uid string)( *SessionObject, error){
	v, ok := m.sessionQueue[uid]
	if ok {
		return v, nil
	} else {
		return nil, errors.New("uid not existing")
	}
}

func (m* SessionManager) Destroy(uid string) error {
	v, ok := m.sessionQueue[uid]
	if ok {
		// 设置所有设备为 online,清除绑定会话
		v.target.LeaveSession()
		v.client.LeaveSession()
		// 从会话列表删除会话
		delete(m.sessionQueue, uid)
		// 从设备列表删除设备
		delete( TargetManagerInstance().targetQueue, v.target.name)
		delete( ClientManagerInstance().targetQueue, v.client.name)
		return nil
	} else {
		return errors.New("uid not existing")
	}
}