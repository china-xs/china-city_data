/**
 @author: xs
 @date: 2021/10/28
 @Description:
**/
package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	mainURL = "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2020/index.html"
)

// 省
func GetSf(get bool, num int) {
	var err error

	example := Client(mainURL)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(example))
	if err != nil {
		log.Fatal(err)
	}
	mainURL := mainURL[:strings.Index(mainURL, "/index.html")]
	dom.Find("a").Each(func(i int, selection *goquery.Selection) {
		t, _ := selection.Html()
		link, _ := selection.Attr("href")
		name := t[:strings.Index(t, "<br/>")]
		if !get {
			save(0, 1, name, "", name, "", i+1)
		}
		if get && (i+1) >= num {
			var district MDistrict
			db.Table("cn_district").Where("name = ?", name).First(&district)
			newId := district.ID
			handler(
				mainURL+"/"+link,
				strconv.Itoa(int(newId)),
				name,
				newId,
				2,
			)
			//fmt.Println(name)
			//if true {
			//	GetSq(
			//		mainURL+"/"+link,
			//		strconv.Itoa(int(newId)),
			//		name,
			//		newId,
			//	)
			//}

		}
	})
	fmt.Println("finish")
}

// 市
func GetSq(url, pIds, merger string, id int32) {
	var level int32 = 2
	example := Client(url)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(example))
	if err != nil {
		log.Fatal(err)
	}
	urlPrix := url[:strings.LastIndex(url, "/")]
	var code, name string
	dom.Find("a").Each(func(i int, selection *goquery.Selection) {
		t, _ := selection.Html()
		if i%2 == 0 {
			code = t
		} else {
			name = t
			newId := save(id, level, name, pIds, merger+","+name, code, i+1)
			link, ok := selection.Attr("href")
			if ok {
				GetXs(
					urlPrix+"/"+link,
					pIds+","+strconv.Itoa(int(newId)),
					merger+","+name,
					newId,
				)
			}
		}
	})
	//log.Printf("Go's time.After example:\n%s", example)

}

// 县级市
func GetXs(url, pIds, merger string, id int32) {
	var level int32 = 3
	example := Client(url)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(example))
	if err != nil {
		log.Fatal(err)
	}
	urlPrix := url[:strings.LastIndex(url, "/")]
	var code, name string
	dom.Find("a").Each(func(i int, selection *goquery.Selection) {
		t, _ := selection.Html()
		if i%2 == 0 {
			code = t
		} else {
			name = t
			newId := save(id, level, name, pIds, merger+","+name, code, i+1)
			link, ok := selection.Attr("href")
			if ok {
				GetJd(
					urlPrix+"/"+link,
					pIds+","+strconv.Itoa(int(newId)),
					merger+","+name,
					newId,
				)
			}
		}
	})
}

// 街道
func GetJd(url, pIds, merger string, id int32) {
	var level int32 = 4
	example := Client(url)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(example))
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(merger, "东莞市") || strings.Contains(merger, "中山市") {
		dom.Find(".villagetr").Each(func(i int, selection *goquery.Selection) {
			t, _ := selection.Html()
			s1 := strings.Index(t, ">")
			s2 := strings.Index(t, "</td>")
			s3 := strings.LastIndex(t, "<td>")
			s4 := strings.LastIndex(t, "</td>")
			code := t[s1+1 : s2]
			name := t[s3+4 : s4]
			save(id, level, name, pIds, merger+","+name, code, i+1)
		})
	} else {
		urlPrix := url[:strings.LastIndex(url, "/")]
		var code, name string

		dom.Find("a").Each(func(i int, selection *goquery.Selection) {
			t, _ := selection.Html()
			if i%2 == 0 {
				code = t
			} else {
				name = t
				newId := save(id, level, name, pIds, merger+","+name, code, i+1)
				link, ok := selection.Attr("href")
				if ok {
					GetJw(
						urlPrix+"/"+link,
						pIds+","+strconv.Itoa(int(newId)),
						merger+","+name,
						newId,
					)
				}
			}
		})
	}
}

// 居委会
func GetJw(url, pIds, merger string, id int32) {
	var level int32 = 5
	example := Client(url)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(example))
	if err != nil {
		log.Fatal(err)
	}
	dom.HasClass("")
	dom.Find(".villagetr").Each(func(i int, selection *goquery.Selection) {
		t, _ := selection.Html()
		s1 := strings.Index(t, ">")
		s2 := strings.Index(t, "</td>")
		s3 := strings.LastIndex(t, "<td>")
		s4 := strings.LastIndex(t, "</td>")
		code := t[s1+1 : s2]
		name := t[s3+4 : s4]
		save(id, level, name, pIds, merger+","+name, code, i+1)
	})
}

type MDistrict struct {
	ID        int32     `gorm:"column:id;type:int;primaryKey" json:"id"`
	PID       int32     `gorm:"column:p_id;type:int unsigned;not null;default:0" json:"p_id"` // 上级菜单id
	PIds      string    `gorm:"column:p_ids;type:varchar(128);not null" json:"p_ids"`
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Level     int32     `gorm:"column:level;type:tinyint(4);not null" json:"level"`
	Merger    string    `gorm:"column:merger;type:varchar(255);not null" json:"merger"`
	Code      string    `gorm:"column:code;type:varchar(255);not null" json:"code"`
	OrderNum  int32     `gorm:"column:order_num;type:int unsigned;not null;default:0" json:"order_num"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

//
//  save
//  @Description:
//  @param pid
//  @param level
//  @param name
//  @param pIds
//  @param merger
//
func Client(url string) (table string) {
	if Cache.IsExist(url) {
		c := Cache.Get(url)
		return c.(string)
	}
	s := rand.Intn(3)
	time.Sleep(time.Second * time.Duration(s))
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	chromedp.Run(ctx,
		chromedp.Navigate(url),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`#Map`),
		chromedp.InnerHTML(`body > table:nth-child(3) > tbody > tr:nth-child(1) > td > table > tbody`, &table),
	)
	if table != "" {
		fmt.Println(url)
		Cache.Put(url, table, expire)
	}
	return
}
