package models

import (
	"errors"
	"fmt"
	"reflect"
	"server/common"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Script struct {
	CreateAt time.Time `orm:"column(createAt);type(timestamp);null" description:"添加时间"`
	Creator  string    `orm:"column(creator)" description:"添加人id"`
	Desc     string    `orm:"column(desc);size(512);null" description:"脚本说明"`
	Id       int       `orm:"column(id);auto"`
	Name     string    `orm:"column(name);size(128)" description:"脚本名称"`
	Path     string    `orm:"column(path);size(512)" description:"脚本储存位置"`
	UpdateAt time.Time `orm:"column(updateAt);type(timestamp);null" description:"更新时间"`
}

func (t *Script) TableName() string {
	return "script_library"
}

func init() {
	common.GetOrm()
	orm.RegisterModel(new(Script))
}

// AddScriptLibrary insert a new Script into database and returns
// last inserted Id on success.
func AddScriptLibrary(m *Script) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetScriptLibraryById retrieves Script by Id. Returns error if
// Id doesn't exist
func GetScriptLibraryById(id int) (v *Script, err error) {
	o := orm.NewOrm()
	v = &Script{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllScriptLibrary retrieves all Script matches certain condition. Returns empty list if
// no records exist
func GetAllScriptLibrary(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Script))
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

	var l []Script
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

// UpdateScriptLibrary updates Script by Id and returns error if
// the record to be updated doesn't exist
func UpdateScriptLibraryById(m *Script) (err error) {
	o := orm.NewOrm()
	v := Script{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteScriptLibrary deletes Script by Id and returns error if
// the record to be deleted doesn't exist
func DeleteScriptLibrary(id int) (err error) {
	o := orm.NewOrm()
	v := Script{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		ossName := v.Path
		if num, err = o.Delete(&Script{Id: id}); err == nil {
			if len(ossName) > 0 {
				common.MinioManagerInstance().Delete(ossName)
			}
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
