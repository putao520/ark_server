package trans

import (
	"errors"
	"github.com/gorilla/websocket"
)

type DeviceManager struct{
	targetQueue map[string]*DeviceObject
}

func DeviceManagerNew() *DeviceManager{
	return &DeviceManager{
		targetQueue: make(map[string]*DeviceObject),
	}
}

func (m *DeviceManager) Join(name string, client *websocket.Conn) *DeviceObject{
	 do, ok := m.targetQueue[name]
	 if ok {
		 return do
	 } else {
		 do = DeviceObjectNew(name, client)
		 m.targetQueue[name] = do
	 }
	 return do
}

func (m *DeviceManager) Leave(name string){
	do, ok := m.targetQueue[name]
	if ok {
		delete(m.targetQueue, name)
	}
	do.LeaveSession()
}

func (m *DeviceManager) Device(name string) (*DeviceObject, error){
	dev, ok := m.targetQueue[name]
	if ok {
		return dev, nil
	} else{
		return nil, errors.New("device is not existing")
	}
}