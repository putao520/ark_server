package common

import (
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

var inited = false

func GetOrm() {
	if !inited {
		inited = true
		sqlconn, err := beego.AppConfig.String("sqlconn")
		if err != nil {
			panic(err)
		}
		database, err := beego.AppConfig.String("database")
		if err != nil {
			panic(err)
		}
		orm.RegisterDriver(database, orm.DRMySQL)
		orm.RegisterDataBase("default", database, sqlconn)
	}

}
