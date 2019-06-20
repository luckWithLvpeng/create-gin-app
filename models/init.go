package models

import (
	"flag"
	"fmt"
	"log"

	"eme/pkg/config"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	// DB 数据库链接池
	db  *gorm.DB
	err error
	// 是否初始化数据库
	syncdb bool
)

func init() {
	flag.BoolVar(&syncdb, "syncdb", false, "bool 是否清空数据库 默认: false")
	flag.Parse()
	initDB()
	creatTable()
}

func initDB() {
	sec := config.DefaultConfig.Section("database")
	user := sec.Key("user").MustString("root")
	password := sec.Key("password").MustString("q1w2e3r4")
	host := sec.Key("host").MustString("0.0.0.0:3306")
	database := sec.Key("database").MustString("emedb")
	url := fmt.Sprintf("%s:%s@(%s)/?charset=utf8&parseTime=True&loc=Local", user, password, host)
	db, err = gorm.Open("mysql", url)

	if err != nil {
		log.Fatalf("connect database %s error: %v\n", url, err)
	}
	if syncdb {
		db = db.Exec(fmt.Sprintf("drop database %s", database))
	}
	db = db.Exec(fmt.Sprintf("create database if not exists %s", database))
	db = db.Exec(fmt.Sprintf("use %s;", database))
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func creatTable() {
	if !db.HasTable(&Role{}) {
		db.CreateTable(&Role{})
	}
	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}
}
