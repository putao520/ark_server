package models

import "server/trans"

type RoomCreateRequest struct {
	DeviceId string
	ClientId string
}

type RoomCreateStatus int

const (
	Success           RoomCreateStatus = 0
	TargetNotExisting RoomCreateStatus = 1
	TargetIsBusy      RoomCreateStatus = 2
	ClientNotExisting RoomCreateStatus = 3
)

// AddRoom 返回会话id或者空字符串
func AddRoom(v RoomCreateRequest) (string, RoomCreateStatus) {
	// 当前 target 已经在某会话,终止创建该会话
	to, err := trans.TargetManagerInstance().Device(v.DeviceId)
	if err != nil {
		return "", TargetNotExisting
	}
	if to.GetStatus() == trans.Busy {
		return "", TargetIsBusy
	}

	// 当前 client 已经在某会话,自动关闭某会话
	co, err := trans.ClientManagerInstance().Device(v.ClientId)
	if err != nil {
		return "", ClientNotExisting
	}
	if co.GetStatus() == trans.Busy {
		se := co.GetSession()
		if se != nil {
			// 断开当前会话目标机连接
			err := se.GetTarget().LeaveSession().GetConn().Close()
			if err != nil {
				return "", 0
			}
			err = trans.SessionManagerInstance().Destroy(se.GetId())
			if err != nil {
				return "", 0
			}
		}
	}

	// 创建会话
	return trans.SessionManagerInstance().Build(to, co).GetId(), Success
}
