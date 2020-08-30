package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"math"
	"reflect"
	. "school-web/common"
	"strings"
	"time"
)

type File struct {
	Id         int       `orm:"column(id);auto"`
	CreateTime time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	Path       string    `orm:"column(path);size(255)" description:"路径"`
	Type       int       `orm:"column(type);null" description:"文件类型"`
}

func (t *File) TableName() string {
	return "file"
}

func init() {
	orm.RegisterModel(new(File))
}

// AddFile insert a new File into database and returns
// last inserted Id on success.
func AddFile(m *File) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

//add files
func AddFiles(length int, m *[]File) (successNumbers int64, err error) {
	o := orm.NewOrm()
	successNumbers, err = o.InsertMulti(length, m)

	return successNumbers, err
}

// GetFileById retrieves File by Id. Returns error if
// Id doesn't exist
func GetFileById(id int) (v *File, err error) {
	o := orm.NewOrm()
	v = &File{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllFile retrieves all File matches certain condition. Returns empty list if
// no records exist
func GetAllFile(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(File))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []File
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// GetAllFile retrieves all File matches certain condition. Returns empty list if
// no records exist
func GetAllFilePage(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (filePage *Page, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(File))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	page := new(Page)
	var ml []interface{}
	var l []File
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {

		count, err := qs.Count()
		if err != nil {
			return nil, errors.New("Error: count error")
		}
		page.Total = count
		page.Pages = int(math.Ceil((float64(count) / float64(limit))))
		page.PageSize = limit
		page.PageNum = offset

		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		page.List = ml
		return page, nil
	}
	return nil, err
}

// UpdateFile updates File by Id and returns error if
// the record to be updated doesn't exist
func UpdateFileById(m *File) (err error) {
	o := orm.NewOrm()
	v := File{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteFile deletes File by Id and returns error if
// the record to be deleted doesn't exist
func DeleteFile(id int) (err error) {
	o := orm.NewOrm()
	v := File{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&File{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//通过id列表查询文件列表，ID用英文逗号分隔，例如：1，2，3
func FindFilesByIds(fileIds string) (files []File, err error) {
	o := orm.NewOrm()
	var fileList []File
	_, err = o.Raw("select * from file where id in (" + fileIds + ")").QueryRows(&fileList)
	if err != nil {
		return nil, err
	}
	return fileList, nil
}
