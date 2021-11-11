package controllers

import (
	"encoding/json"
	"errors"
	"server/jwt"
	"server/models"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// UserController operations for User
type UserController struct {
	beego.Controller
}

// URLMapping ...
func (c *UserController) URLMapping() {
	c.Mapping("Login", c.Login)
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {int} models.User
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Post() {
	var v models.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddUser(&v); err == nil {
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
// @Description get User by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UserController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	v, err := models.GetUserById(idStr)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = err.Error()
	} else {
		v.Password = ""
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.User
// @Failure 403
// @router / [get]
func (c *UserController) GetAll() {
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

	l, err := models.GetAllUser(query, fields, sortby, order, offset, limit)
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
// @Description update the User
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 403 :id is not int
// @router / [put]
func (c *UserController) Put() {
	userInfoStr := c.Ctx.Input.Param("UserInfo")
	userInfo := &models.User{Id: ""}
	err := json.Unmarshal([]byte(userInfoStr), userInfo)
	if err != nil {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "need login"
	}
	idStr := userInfo.Id
	v := models.User{Id: idStr}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateUserById(&v); err == nil {
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
// @Description delete the User
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UserController) Delete() {
	/*
		idStr := c.Ctx.Input.Param(":id")
		if err := models.DeleteUser(idStr); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = err.Error()
		}
	*/
	c.Ctx.Output.SetStatus(405)
	c.Data["json"] = "Not Support Method"
	c.ServeJSON()
}

// Login ...
// @Title Login
// @Description User login
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {string} token
// @Failure 403 body is empty
// @router /Login [post]
func (c *UserController) Login() {
	var loginInfo map[string]string
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &loginInfo); err == nil {
		username, ok := loginInfo["username"]
		if !ok {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = "need username"
		}
		password, ok := loginInfo["password"]
		if !ok {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = "need password"
		}
		v, err := models.GetUserByIdAndPassword(username, password)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = err.Error()
		}
		token, err := jwt.GenerateToken(v, 0)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = err.Error()
		} else {
			c.Data["json"] = token
		}
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = err
	}
	c.ServeJSON()
}
