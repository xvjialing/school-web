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

type Article struct {
	Id         int       `orm:"column(id);auto"`
	Title      string    `orm:"column(title);size(255)"`
	Subtitle   string    `orm:"column(subtitle);size(255);null"`
	Content    string    `orm:"column(content)"`
	CreateTime time.Time `orm:"column(create_time);type(timestamp);auto_now"`
	Author     string    `orm:"column(author);size(255)"`
	Type       string    `orm:"column(type);size(255);null"  description:"文章类型"`
	FileIdList string    `orm:"column(file_id_list);size(255);null"  description:"附件ID列表，英文逗号分隔，例如：1,2,3,4,5"`
}

type ArticleFiles struct {
	Id         int
	Title      string
	Subtitle   string
	Content    string
	CreateTime time.Time
	Author     string
	Type       string
	FileIdList []File
}

func ArticleToArticleFiles(article Article) (articleFile *ArticleFiles, err error) {
	var articleFiles ArticleFiles
	articleFiles.Id = article.Id
	articleFiles.Title = article.Title
	articleFiles.Subtitle = article.Subtitle
	articleFiles.Content = article.Content
	articleFiles.CreateTime = article.CreateTime
	articleFiles.Author = article.Author
	articleFiles.Type = article.Type

	fileIdsStr := article.FileIdList

	if len(fileIdsStr) != 0 {
		files, err := FindFilesByIds(fileIdsStr)
		if err != nil {
			return nil, err
		}
		articleFiles.FileIdList = files
	}
	return &articleFiles, nil
}

func ArticalListToArticleFilsList(articleList []Article) (articleFileList *[]ArticleFiles, err error) {
	var articleFilesLists []ArticleFiles

	for _, v := range articleList {
		articleFiles, err := ArticleToArticleFiles(v)
		if err != nil {
			return nil, err
		}
		articleFilesLists = append(articleFilesLists, *articleFiles)
	}
	return &articleFilesLists, nil
}

func (t *Article) TableName() string {
	return "article"
}

func init() {
	orm.RegisterModel(new(Article))
}

// AddArticle insert a new Article into database and returns
// last inserted Id on success.
func AddArticle(m *Article) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetArticleById retrieves Article by Id. Returns error if
// Id doesn't exist
func GetArticleById(id int) (v *Article, err error) {
	o := orm.NewOrm()
	v = &Article{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllArticle retrieves all Article matches certain condition. Returns empty list if
// no records exist
func GetAllArticle(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Article))
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

	var l []Article
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

// GetAllArticle retrieves all Article matches certain condition. Returns empty list if
// no records exist
func GetAllArticlePage(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (articlePage *Page, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Article))
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
	var l []Article
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

		var articleFileList []interface{}
		for _, v := range ml {
			articleFile, err := ArticleToArticleFiles(v.(Article))
			if err != nil {
				continue
			}
			articleFileList = append(articleFileList, *articleFile)
		}
		page.List = articleFileList
		return page, nil
	}
	return nil, err
}

// UpdateArticle updates Article by Id and returns error if
// the record to be updated doesn't exist
func UpdateArticleById(m *Article) (err error) {
	o := orm.NewOrm()
	v := Article{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteArticle deletes Article by Id and returns error if
// the record to be deleted doesn't exist
func DeleteArticle(id int) (err error) {
	o := orm.NewOrm()
	v := Article{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Article{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
