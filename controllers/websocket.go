package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gorilla/websocket"
	"net/http"
	"server/trans"
	"strconv"
)

type DeviceClass int

const (
	Target DeviceClass = 0
	Client DeviceClass = 1
	LastClass DeviceClass = 2
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024 * 1024 * 10,
	WriteBufferSize: 1024 * 1024 * 10,
}

type WebSocketController struct {
	beego.Controller
}

/*
func (m *WebSocketController) doMessageTarget(conn *websocket.Conn, se *trans.SessionObject) error {
	for {
		t, c, e := conn.ReadMessage()
		if t == -1 || e != nil {
			return e
		}
		e = se.ForwardToClient(t, c)
		if e != nil {
			return e
		}
	}
	return nil
}


func (m *WebSocketController) doMessageClient(conn *websocket.Conn, se *trans.SessionObject) error {
	for {
		t, c, e := conn.ReadMessage()
		if t == -1 || e != nil {
			return e
		}
		e = se.ForwardToTarget(t, c)
		if e != nil {
			return e
		}
	}
	return nil
}
*/

func (m *WebSocketController) accept() (*websocket.Conn, string, bool){
	param, err := m.Input() // (":deviceId")
	if err != nil {
		http.Error(m.Ctx.ResponseWriter, "need connect param", 400)
		return nil, "", false
	}
	deviceId := param.Get("deviceId")
	if len(deviceId) == 0 {
		http.Error(m.Ctx.ResponseWriter, "need deviceId", 400)
		return nil, "", false
	}
	ws, err :=  upgrader.Upgrade(m.Ctx.ResponseWriter, m.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(m.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return nil, "", false
	} else if err != nil {
		fmt.Println(ws)
		return nil, "", false
	}
	return ws, deviceId, true
}

func (m *WebSocketController) OnConnect() {
	conn, deviceId, ok := m.accept()
	if !ok {
		return
	}
	param, _ := m.Input()
	deviceClass := param.Get("deviceClass")
	if len(deviceClass) == 0 {
		http.Error(m.Ctx.ResponseWriter, "need deviceClass", 400)
		conn.Close()
		return
	}
	dcv, err := strconv.Atoi(deviceClass)
	if err != nil || dcv >= int(LastClass) {
		http.Error(m.Ctx.ResponseWriter, "deviceClass invalid", 400)
		conn.Close()
		return
	}
	// 添加到对应角色管理器
	switch dcv {
	case int(Target):
		do := trans.TargetManagerInstance().Join(deviceId, conn)
		dse := do.GetSession()
		if dse != nil {
			dse.RestoreTarget(conn)
		}
		break
	case int(Client):
		do := trans.ClientManagerInstance().Join(deviceId, conn)
		dse := do.GetSession()
		if dse != nil {
			dse.RestoreClient(conn)
		}
		break
	default:
	}

	// 中断连接
	conn.SetCloseHandler( func(code int , text string) error {
		var se *trans.SessionObject = nil
		switch dcv {
		case int(Target):
			targetManager := trans.TargetManagerInstance()
			deviceObject, err := targetManager.Device(deviceId)
			if err == nil {
				deviceObject.SetStatus(trans.Lost)
			}
			se = deviceObject.GetSession()
			/*
			// 单个对象,断开连接时,不自动离开会话

			targetManager.Leave(deviceId)
			*/
			break
		case int(Client):
			clientManager := trans.ClientManagerInstance()
			deviceObject, err := clientManager.Device(deviceId)
			if err == nil {
				deviceObject.SetStatus(trans.Lost)
			}
			se = deviceObject.GetSession()
			/*
			// 单个对象,断开连接时,不自动离开会话
			clientManager.Leave(deviceId)
			*/
			break
		default:
		}
		// 如果退出的conn包含会话,而且会话内所有端都离线了,那么自动撤销会话
		if se != nil && se.CanDestroy() {
			trans.SessionManagerInstance().Destroy(se.GetId())
		}
		return nil
	})
}