package main

import (
	"fmt"
	"github.com/huntsman-li/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"path"
	"time"
)

var (
	db     *gorm.DB
	Cache  cache.Cache
	expire int64 = 15768000
)

const (
	host   = "192.168.56.12"
	user   = "root"
	pwd    = "root"
	port   = 3306
	dbname = ""
	preFix = "oa_"
)

func main() {
	//sk := Client("https://baidu.com")
	//fmt.Printf("sk:%v,skType:%T\n ss:%v",sk,sk, sk=="")
	GetSf(true, 1)

}

func init() {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user,
		pwd,
		host,
		port,
		dbname,
	)
	config := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   preFix,
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}
	db, _ = gorm.Open(mysql.Open(dsn), &config)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(6 * time.Hour)

	// set 缓存
	pwd, _ := os.Getwd()
	dir := path.Join(pwd, "caches")
	cache, err := cache.Cacher(cache.Options{
		Adapter:       "file",
		AdapterConfig: dir,
		Interval:      2,
	})
	if err != nil {
		panic(err)
	}
	Cache = cache

}
