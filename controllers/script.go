package controllers

import (
	"encoding/json"
	"errors"
	"server/common"
	"server/jwt"
	"server/models"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// ScriptController operations for Script
type ScriptController struct {
	beego.Controller
}

// URLMapping ...
func (c *ScriptController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Upload", c.Upload)
	c.Mapping("GeneratorOnceUrl", c.GeneratorOnceUrl)
}

// Post ...
// @Title Post
// @Description create Script
// @Param	body		body 	models.Script	true		"body for Script content"
// @Success 201 {int} models.Script
// @Failure 403 body is empty
// @router / [post]
func (c *ScriptController) Post() {
	var v models.Script
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddScriptLibrary(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = err.Error()
		}
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Script by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Script
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ScriptController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetScriptLibraryById(id)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Script
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Script
// @Failure 403
// @router / [get]
func (c *ScriptController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllScriptLibrary(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Script
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Script	true		"body for Script content"
// @Success 200 {object} models.Script
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ScriptController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Script{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateScriptLibraryById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = err.Error()
		}
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Script
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ScriptController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteScriptLibrary(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Upload ...
// @Title Upload
// @Description Upload the Script
// @Success 200 {string} Upload success!
// @Failure 403 id is empty
// @router /Upload/ [post]
func (c *ScriptController) Upload() {
	// 获得当前会话信息
	userInfoStr := c.Ctx.Input.Param("UserInfo")
	userInfo, err := jwt.GetSessionInfo([]byte(userInfoStr))
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "need login"
		c.ServeJSON()
		return
	}

	f, h, err := c.GetFile("file")
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	defer f.Close()
	buffer := make([]byte, h.Size)
	_, err = f.Read(buffer)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	// 上传文件到 oss
	ossFileUrl, err := common.MinioManagerInstance().Upload(h.Filename, buffer)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}

	desc := c.GetString("desc")
	// 添加数据到数据库
	logsFileInfo := &models.Script{
		CreateAt: time.Time{},
		Creator:  userInfo.Id,
		Desc:     desc,
		Name:     h.Filename,
		Path:     ossFileUrl,
		UpdateAt: time.Time{},
	}
	// 添加文件记录到数据库
	if _, err := models.AddScriptLibrary(logsFileInfo); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = logsFileInfo
	} else {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GeneratorOnceUrl ...
// @Title GeneratorOnceUrl
// @Description generator once download file url
// @Param	id		objectName 	string	true		"The object name you want to download"
// @Success 200 {string} once download url!
// @Failure 403 id is empty
// @router /GeneratorOnceUrl/:id [get]
func (c *ScriptController) GeneratorOnceUrl() {
	objectName := c.Ctx.Input.Param(":id")
	if len(objectName) == 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "need file name"
		c.ServeJSON()
		return
	}
	downloadUrl, err := common.MinioManagerInstance().PresignedGet(objectName, "30s")
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = downloadUrl
	}
	c.ServeJSON()
}
