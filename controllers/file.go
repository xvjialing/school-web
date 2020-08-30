package controllers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"os"
	"school-web/common"
	"school-web/models"
	"strconv"
	"strings"
	"time"
)

// 文件相关操作
type FileController struct {
	BaseController
}

// URLMapping ...
func (c *FileController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create File
// @Param	access_token	header	string	true	"access_token"
// @Param	file		form 	file	true		"body for File content"
// @Success 201 {int} models.File
// @Failure 403 body is empty
// @router / [post]
func (c *FileController) Post() {

	var fileList []models.File

	files, err := c.GetFiles("file")
	if err != nil {
		c.Data["json"] = common.Failed(400, err.Error())
		c.ServeJSON()
	}

	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			c.Data["json"] = common.Failed(400, err.Error())
			c.ServeJSON()
			return
		}
		fileName := uuid.New().String() + "_" + files[i].Filename

		//创建目录
		uploadDir := "static/upload/" + time.Now().Format("2006/01/02/")
		err = os.MkdirAll(uploadDir, 0777)
		if err != nil {
			c.Data["json"] = common.Failed(400, err.Error())
			c.ServeJSON()
			return
		}

		path := uploadDir + fileName

		dst, err := os.Create(path)
		defer dst.Close()
		if err != nil {
			c.Data["json"] = common.Failed(400, err.Error())
			c.ServeJSON()
			return
		}
		_, err = io.Copy(dst, file)
		if err != nil {
			c.Data["json"] = common.Failed(400, err.Error())
			c.ServeJSON()
			return
		}
		modelFile := models.File{
			CreateTime: time.Now(),
			Path:       path,
			Type:       1,
		}
		fileList = append(fileList, modelFile)

	}

	_, err = models.AddFiles(len(fileList), &fileList)

	if err != nil {
		c.Data["json"] = common.Failed(400, err.Error())
		c.ServeJSON()
		return
	}

	c.Data["json"] = common.Succes(fileList)

	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get File by id
// @Param	access_token	header	string	true	"access_token"
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.File
// @Failure 403 :id is empty
// @router /:id [get]
func (c *FileController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetFileById(id)
	if err != nil {
		c.Data["json"] = common.Failed(400, err.Error())
	} else {
		c.Data["json"] = common.Succes(v)
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get File
// @Param	access_token	header	string	true	"access_token"
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	pageSize	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	pageNum	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.File
// @Failure 403
// @router / [get]
func (c *FileController) GetAll() {
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

	l, err := models.GetAllFilePage(query, fields, sortby, order, (currentPage-1)*pageSize, pageSize)
	if err != nil {
		c.Data["json"] = common.Failed(400, err.Error())
	} else {
		c.Data["json"] = common.Succes(l)
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the File
// @Param	access_token	header	string	true	"access_token"
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.File	true		"body for File content"
// @Success 200 {object} models.File
// @Failure 403 :id is not int
// @router /:id [put]
func (c *FileController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.File{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateFileById(&v); err == nil {
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
// @Description delete the File
// @Param	access_token	header	string	true	"access_token"
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *FileController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteFile(id); err == nil {
		c.Data["json"] = common.Succes("OK")
	} else {
		c.Data["json"] = common.Failed(400, err.Error())
	}
	c.ServeJSON()
}
