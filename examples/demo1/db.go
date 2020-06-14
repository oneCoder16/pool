package main

import (
	"fmt"

	// MySQL driver.
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Options struct {
	UserName string
	Psw      string
	Addr     string
	Name     string
}

func OpenDB(opts *Options) (*gorm.DB, error) {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		opts.UserName,
		opts.Psw,
		opts.Addr,
		opts.Name,
		true,
		//"Asia/Shanghai"),
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		return nil, err
	}

	setupDB(db)

	return db, nil
}

func setupDB(db *gorm.DB) {
	db.LogMode(true)
	//db.DB().SetMaxOpenConns(20000) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(0) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}
