package main

import (
	"fmt"
	"github.com/huntsman-li/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	db     *gorm.DB
	Cache  cache.Cache
	expire int64 = 15768000
)

const (
	host   = "192.168.56.12"
	user   = "db_user"
	pwd    = "db_pass"
	port   = 3306
	dbname = "yjf_scrm_1000001"
	preFix = "oa_"
	table  = "m_district"
)

func main() {
	//sk := Client("https://baidu.com")
	//fmt.Printf("sk:%v,skType:%T\n ss:%v",sk,sk, sk=="")
	//GetSf(true, 1)
	var results []MDistrict
	var pids string
	url := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/index.html"
	for i := 1; i <= 5; i++ {
		if i > 1 {
			//db.Table(table).Create(&save)
			db.Table(table).Where("level = ?", i-1).FindInBatches(&results, 10, func(tx *gorm.DB, batch int) error {
				for _, result := range results {
					// 批量处理找到的记录
					ck := strconv.Itoa(int(result.ID))
					if Cache.IsExist(ck) {
						c := Cache.Get(ck)
						url = c.(string)
						if result.PIds != "" {
							pids = result.PIds + "," + ck
						} else {
							pids = ck
						}
						merger := result.Merger
						handlerV2(url, pids, merger, result.ID, int32(i))
					}
				}
				return nil
			})
		} else {
			handlerV2(url, "", "", 0, 1)
		}
	}
	/*
		url := "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/44.html"
		url ="http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/44/20/442000002.html"
		html := GetHtml(url)
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			log.Fatal(err)
		}
		//urlPrix := url[:strings.LastIndex(url, "/")]
		//var code, name string
		hasALink := false
		dom.Find("a").Each(func(i int, selection *goquery.Selection) {
			hasALink = true

		})
		if !hasALink {
			html := GetHtml(url)
			dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
			if err != nil {
				log.Fatal(err)
			}
			//urlPrix := url[:strings.LastIndex(url,"/")]
			dom.Find(".villagetr").Each(func(i int, selection *goquery.Selection) {
				t, _ := selection.Html()
				s1 := strings.Index(t, ">")
				s2 := strings.Index(t, "</td>")
				s3 := strings.LastIndex(t, "<td>")
				s4 := strings.LastIndex(t, "</td>")
				code := t[s1+1 : s2]
				name := t[s3+4 : s4]
				fmt.Printf("code:%v,name:%v\n", code,name)
				//save(id, level, name, pIds, merger+","+name, code, i+1)
			})
		}

	*/

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
