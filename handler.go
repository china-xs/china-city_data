/**
 @author: xs
 @date: 2021/11/3
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

func handlerV2(url, pids, merger string, id, level int32) {
	html := GetHtml(url)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}
	urlprix := url[:strings.LastIndex(url, "/")]
	var code, name string
	hasLink := false
	var j = 1
	dom.Find("a").Each(func(i int, selection *goquery.Selection) {
		hasLink = true
		t, _ := selection.Html()
		if level == 1 {
			name = t[:strings.Index(t, "<br/>")]
			newId := save(id, level, name, pids, name, code, j)
			j++
			link, _ := selection.Attr("href")
			ck := strconv.Itoa(int(newId))
			cv := urlprix + "/" + link
			Cache.Put(ck, cv, expire)
		} else {
			if i%2 == 0 {
				code = t
			} else {
				name = t
				newId := save(id, level, name, pids, merger+","+name, code, j)
				link, _ := selection.Attr("href")
				//cv := urlprix+"/"+link,
				ck := strconv.Itoa(int(newId))
				cv := urlprix + "/" + link
				Cache.Put(ck, cv, expire)
				j++
			}
		}

	})
	if !hasLink {
		dom.Find(".villagetr").Each(func(i int, selection *goquery.Selection) {
			t, _ := selection.Html()
			s1 := strings.Index(t, ">")
			s2 := strings.Index(t, "</td>")
			s3 := strings.LastIndex(t, "<td>")
			s4 := strings.LastIndex(t, "</td>")
			code := t[s1+1 : s2]
			name := t[s3+4 : s4]
			save(id, level, name, pids, merger+","+name, code, i+1)
		})
	}
}

func handler(url, pIds, merger string, id, level int32) {
	html := GetHtml(url)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}
	urlPrix := url[:strings.LastIndex(url, "/")]
	var code, name string
	hasLink := false
	dom.Find("a").Each(func(i int, selection *goquery.Selection) {
		hasLink = true
		t, _ := selection.Html()
		if i%2 == 0 {
			code = t
		} else {
			name = t
			newId := save(id, level, name, pIds, merger+","+name, code, i+1)
			link, ok := selection.Attr("href")
			if ok {
				if level == 4 || (level == 3 &&
					strings.Contains(merger, "东莞市") ||
					strings.Contains(merger, "中山市")) {
					handler1(
						urlPrix+"/"+link,
						pIds+","+strconv.Itoa(int(newId)),
						merger+","+name,
						newId,
						level+1,
					)
				} else {
					handler(
						urlPrix+"/"+link,
						pIds+","+strconv.Itoa(int(newId)),
						merger+","+name,
						newId,
						level+1,
					)
				}

			}
		}
	})
	if !hasLink {
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
}

// url pids merger pid level
//url, pIds, merger string, id, level int32

//东莞、中山 居委会
func handler1(url, pIds, merger string, id, level int32) {
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
		save(id, level, name, pIds, merger+","+name, code, i+1)
	})
}

func GetHtml(url string) (html string) {
	if Cache.IsExist(url) {
		c := Cache.Get(url)
		Cache.Put(url, c, expire)
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
		chromedp.InnerHTML(`body > table:nth-child(3) > tbody > tr:nth-child(1) > td > table > tbody`, &html),
	)
	if html != "" {
		fmt.Println(url)
		Cache.Put(url, html, expire)
	}
	return
}

func save(pid, level int32, name, pIds, merger, code string, orderNun int) int32 {
	//type
	save := MDistrict{
		PID:      pid,
		PIds:     pIds,
		Level:    level,
		Name:     name,
		Merger:   merger,
		Code:     code,
		OrderNum: int32(orderNun),
	}
	db.Table("m_district").Create(&save)
	return save.ID
}
