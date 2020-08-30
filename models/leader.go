package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"log"
	"math"
	"reflect"
	. "school-web/common"
	"strconv"
	"strings"
)

type Leader struct {
	Id        int    `form:"_" orm:"column(id);auto"`
	Name      string `orm:"column(name);size(255)" description:"姓名"`
	AvaterUrl string `orm:"column(avater_url);size(255);null" description:"头像路径"`
	Detail    string `orm:"column(detail);null" description:"详情"`
	Title     string `orm:"column(title);size(255);null" description:"标题"`
	Index     int64  `form:"_" orm:"column(index);null"  description:"排序序号"`
}

func (t *Leader) TableName() string {
	return "leader"
}

func init() {
	orm.RegisterModel(new(Leader))
}

// AddLeader insert a new Leader into database and returns
// last inserted Id on success.
func AddLeader(m *Leader) (id int64, err error) {
	m.Index = Count() + 1
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetLeaderById retrieves Leader by Id. Returns error if
// Id doesn't exist
func GetLeaderById(id int) (v *Leader, err error) {
	o := orm.NewOrm()
	v = &Leader{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLeader retrieves all Leader matches certain condition. Returns empty list if
// no records exist
func GetAllLeader(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Leader))
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

	var l []Leader
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

// GetAllLeader retrieves all Leader matches certain condition. Returns empty list if
// no records exist
func GetAllLeaderPage(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (leaderPage *Page, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Leader))
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
	var l []Leader
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

// UpdateLeader updates Leader by Id and returns error if
// the record to be updated doesn't exist
func UpdateLeaderById(m *Leader) (err error) {
	o := orm.NewOrm()
	v := Leader{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLeader deletes Leader by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLeader(id int) (err error) {
	o := orm.NewOrm()
	v := Leader{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		index := v.Index
		AllMinusOne(index)
		var num int64
		if num, err = o.Delete(&Leader{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func LeaderIndexPlusOne(id int) (err error) {
	count := Count()
	o := orm.NewOrm()

	leader := Leader{Id: id}
	if o.Read(&leader) == nil {
		if leader.Index == count {
			return errors.New("已经是最大序号")
		}

		changeLeader := Leader{Index: leader.Index + 1}
		err := o.Read(&changeLeader, "Index")
		if err != nil {
			return err
		}
		changeLeader.Index = changeLeader.Index - 1
		update, err := o.Update(&changeLeader, "Index")
		if err != nil || update < 1 {
			return errors.New("update fialed")
		}

		leader.Index = leader.Index + 1
		i, err := o.Update(&leader, "Index")
		if err != nil || i < 1 {
			return errors.New("update fialed")
		}
	}
	return nil
}

func LeaderIndexMinusOne(id int) (err error) {
	o := orm.NewOrm()
	leader := Leader{Id: id}
	if o.Read(&leader) == nil {
		if leader.Index == 1 {
			log.Println("已经是最小序号")
			return errors.New("已经是最小序号")
		}

		changeLeader := Leader{Index: leader.Index - 1}
		err := o.Read(&changeLeader, "Index")
		if err != nil {
			log.Println(err)
			return err
		}
		changeLeader.Index = changeLeader.Index + 1
		update, err := o.Update(&changeLeader, "Index")
		if err != nil || update < 1 {
			log.Println("update fialed")
			return errors.New("update fialed")
		}

		leader.Index = leader.Index - 1
		i, err := o.Update(&leader, "Index")
		if err != nil || i < 1 {
			log.Println("update fialed")
			return errors.New("update fialed")
		}
	}
	return nil
}

func Count() (count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Leader))
	i, _ := qs.Count()
	return i
}

func AllMinusOne(index int64) error {
	o := orm.NewOrm()
	result, e := o.Raw("update leader set `index` = `index` - 1 where `index` > " + strconv.Itoa(int(index))).Exec()
	if e != nil {
		return e
	}
	log.Println(result)
	return nil
}
