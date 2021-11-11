package controllers

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"server/common"
	"server/models"
	"server/trans"
)

type RoomController struct {
	beego.Controller
}

// URLMapping ...
func (u *RoomController) URLMapping() {
	u.Mapping("Post", u.Post)
	u.Mapping("GetAll", u.GetAll)
}

// @Title CreateRoom
// @Description create room
// @Param	body		body 	models.RoomCreateRequest	true		"body for room create"
// @Success 200 {int} models.Session.Id
// @Failure 403 body is empty
// @router / [post]
func (u *RoomController) Post() {
	out := u.Ctx.ResponseWriter
	var roomCreate models.RoomCreateRequest
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &roomCreate)
	if err != nil {
		out.Write(common.ResponseResultNew(trans.Failed, "invalid request", ""))
	}
	uid, status := models.AddRoom(roomCreate)
	if status != models.Success {
		var msg string
		switch status {
		case models.TargetNotExisting:
			msg = "TargetNotExisting"
			break
		case models.ClientNotExisting:
			msg = "ClientNotExisting"
			break
		case models.TargetIsBusy:
			msg = "TargetIsBusy"
			break
		default:
			msg = "Unknown Error"
		}
		out.Write(common.ResponseResultNew(trans.Failed, msg, ""))
	} else {
		out.Write(common.ResponseResultNew(trans.Success, "", uid))
	}
}

// GetAll ...
// @Title Get All
// @Description get Device
// @Success 200
// @Failure 403
// @router / [get]
func (u *RoomController) GetAll() {
	deviceArr := trans.TargetManagerInstance().GetAll()
	result, err := json.Marshal(deviceArr)
	if err != nil {
		u.Ctx.Output.SetStatus(500)
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = string(result)
	}
	u.ServeJSON()
}
