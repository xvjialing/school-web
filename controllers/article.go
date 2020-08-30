package controllers

import (
	"encoding/json"
	"errors"
	"school-web/common"
	"school-web/models"
	"strconv"
	"strings"
)

// 文章相关操作
type ArticleController struct {
	BaseController
}

// URLMapping ...
func (c *ArticleController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Article
// @Param	access_token	header	string	true	"access_token"
// @Param	body		body 	models.Article	true		"body for Article content"
// @Success 201 {int} models.Article
// @Failure 403 body is empty
// @router / [post]
func (c *ArticleController) Post() {
	var v models.Article

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if id, err := models.AddArticle(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			v.Id = int(id)
			articleFiles, err := models.ArticleToArticleFiles(v)
			if err != nil {
				c.Data["json"] = common.Failed(400, err.Error())
			} else {
				c.Data["json"] = common.Succes(articleFiles)
			}

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
// @Description get Article by id
// @Param	access_token	header	string	true	"access_token"
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Article
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ArticleController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetArticleById(id)
	if err != nil {
		c.Data["json"] = common.Failed(400, err.Error())
	} else {
		c.Data["json"] = common.Succes(v)
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Article
// @Param	access_token	header	string	true	"access_token"
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Article
// @Failure 403
// @router / [get]
func (c *ArticleController) GetAll() {
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
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllArticlePage(query, fields, sortby, order, (currentPage-1)*pageSize, pageSize)
	if err != nil {
		c.Data["json"] = common.Failed(400, err.Error())
	} else {
		c.Data["json"] = common.Succes(l)
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Article
// @Param	access_token	header	string	true	"access_token"
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Article	true		"body for Article content"
// @Success 200 {object} models.Article
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ArticleController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Article{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateArticleById(&v); err == nil {
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
// @Description delete the Article
// @Param	access_token	header	string	true	"access_token"
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ArticleController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteArticle(id); err == nil {
		c.Data["json"] = common.Succes("OK")
	} else {
		c.Data["json"] = common.Failed(400, err.Error())
	}
	c.ServeJSON()
}
