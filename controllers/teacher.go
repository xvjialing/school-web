package controllers

import (
	"encoding/json"
	"school-web/common"
	"school-web/models"
	"strconv"
	"strings"
)

// 教师相关操作
type TeacherController struct {
	BaseController
}

// URLMapping ...
func (c *TeacherController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Teacher
// @Param	access_token	header	string	true	"access_token"
// @Param	body		body 	models.Teacher	true		"body for Teacher content"
// @Success 201 {int} models.Teacher
// @Failure 403 body is empty
// @router / [post]
func (c *TeacherController) Post() {
	var v models.Teacher
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddTeacher(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = common.Succes(v)
		} else {
			c.Data["json"] = common.Failed(400, err.Error())
		}
	} else {
		c.Data["json"] = common.Failed(400, err.Error())
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Teacher by id
// @Param	access_token	header	string	true	"access_token"
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Teacher
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TeacherController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetTeacherById(id)
	if err != nil {
		c.Data["json"] = common.Failed(400, err.Error())
	} else {
		c.Data["json"] = common.Succes(v)
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Teacher
// @Param	access_token	header	string	true	"access_token"
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	pageSize	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	pageNum	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Teacher
// @Failure 403
// @router / [get]
func (c *TeacherController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var pageSize int64 = 10
	var currentPage int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// pageSize: 10 (default is 10)
	if v, err := c.GetInt64("pageSize"); err == nil {
		pageSize = v
	}
	// currentPage: 0 (default is 0)
	if v, err := c.GetInt64("pageNum"); err == nil {
		currentPage = v
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
				c.Data["json"] = common.Failed(400, "Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}
	l, err := models.GetAllTeacherPage(query, fields, sortby, order, (currentPage-1)*pageSize, pageSize)

	if err != nil {
		c.Data["json"] = common.Failed(400, err.Error())
	} else {
		c.Data["json"] = common.Succes(l)
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Teacher
// @Param	access_token	header	string	true	"access_token"
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Teacher	true		"body for Teacher content"
// @Success 200 {object} models.Teacher
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TeacherController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Teacher{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateTeacherById(&v); err == nil {
			c.Data["json"] = common.Succes("OK")
		} else {
			c.Data["json"] = common.Failed(400, err.Error())
		}
	} else {
		c.Data["json"] = common.Failed(400, err.Error())
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Teacher
// @Param	access_token	header	string	true	"access_token"
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TeacherController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteTeacher(id); err == nil {
		c.Data["json"] = common.Succes("OK")
	} else {
		c.Data["json"] = common.Failed(400, err.Error())
	}
	c.ServeJSON()
}
