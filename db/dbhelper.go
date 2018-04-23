package db

import (
	"errors"
	"github.com/boxproject/agent/config"
	l4g "github.com/alecthomas/log4go"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/boxproject/agent/model"

)

//数据库处理
type dbHandler struct {
	alisaName string
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
}

//单例db
var db *dbHandler

func InitRDB(cfg config.DataSource) error {

	l4g.Info("DB:Init start")
	defer l4g.Info("DB:Init finish")
	if db == nil {
		db = &dbHandler{alisaName: cfg.AliasName}
	} else {
		l4g.Error("DB:Init fail:already init")
		return errors.New("DB:Init fail:already init")
	}
	// 参数(可选)  设置最大空闲连接
	// 参数(可选)  设置最大数据库连接 (go >= 1.2)
	if err := orm.RegisterDataBase(cfg.AliasName, cfg.DriverName, cfg.Url, cfg.MaxIdle, cfg.MaxConn); err != nil {
		l4g.Error("conn db error:%v", err)
		return errors.New("conn db error")
	}
	orm.Debug = cfg.Debug
	orm.RunSyncdb("default", false, true)
	return nil
}

//获取ormer
func GetNewOrmer(dbName string) orm.Ormer {
	if dbName != "" {
		if db == nil {
			l4g.Error("DB no init")
			return nil
		}
		ormer := orm.NewOrm()
		ormer.Using(dbName)
		return ormer
	} else {
		return GetDefaultNewOrmer()
	}
}

//获取默认ormer
func GetDefaultNewOrmer() orm.Ormer {
	if db == nil {
		l4g.Error("DB no init")
		return nil
	}
	ormer := orm.NewOrm()
	ormer.Using(db.alisaName)
	return ormer
}
