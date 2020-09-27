package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"database/sql"

	htmlquery "github.com/antchfx/xquery/html"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("start qichezhijia_spider: ")
	fmt.Println("--------------------------")
	objurl := "https://car.autohome.com.cn/price/tag-25-0-1.html"
	htmlUTF8 := ConvertToString(spider(objurl), "gbk", "utf-8")
	// fmt.Println(htmlUTF8)

	xpath := "/html/body/div[1]/script[2]"
	parsedInfo := XpathParseHtml(htmlUTF8, xpath)
	jsonListstr := parsedInfo[16 : len(parsedInfo)-1]
	//fmt.Println(jsonListstr)

	// 解析json
	infos := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonListstr), &infos)
	if err != nil {
		// handle with error here
		log.Fatalf("json parse false: %v", err)
	}

	// 链接mysql
	db, err := sql.Open("mysql", "root:rootroot@tcp(127.0.0.1:3306)/qichezhijia?charset=utf8")
	if err != nil {
		fmt.Printf("connect mysql fail ! [%s]", err)
	} else {
		fmt.Println("connect to mysql success")
	}

	cityitems := infos["CityItems"].([]interface{}) // type : []interface {}
	for index, citydetails := range cityitems {
		// 0 ,map[I:110100 N:北京 P:Beijing S:110000]
		id := index
		areacode := citydetails.(map[string]interface{})["I"]
		cityname := citydetails.(map[string]interface{})["N"]
		citypinyin := citydetails.(map[string]interface{})["P"]

		insertsql := fmt.Sprintf(`INSERT into businesscities (id,areacode,cityname,citypinyin) values (%v,%v,"%v","%v")`, int(id), areacode, cityname, citypinyin)
		fmt.Println(insertsql)
		_, err := db.Exec(insertsql)
		if err != nil {
			log.Fatalf("insert failed, err:%v\n", err)
		}
		fmt.Printf("%v : 已写入mysql；", citydetails.(map[string]interface{}))
	}
	db.Close()
}

func spider(objurl string) string {
	// Client类型代表HTTP客户端。它的零值（DefaultClient）是一个可用的使用DefaultTransport的客户端。
	client := &http.Client{}
	req, _ := http.NewRequest("GET", objurl, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	// Do方法发送请求，返回HTTP回复。它会遵守客户端clien设置的策略（如重定向、cookie、认证）。
	resp, err := client.Do(req)
	if err != nil {
		log.Panic("failed to request objurl: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return fmt.Sprintf("%s", body)
}

// XpathParseHtml 根据提供的xpath来解析html获取对应内容；但这里获取的数据还需要再加工；目前今兼容单个标签里的值
func XpathParseHtml(htmlUTF8, xpath string) string {
	root, _ := htmlquery.Parse(strings.NewReader(htmlUTF8))
	tmp := htmlquery.Find(root, xpath)
	return htmlquery.InnerText(tmp[0])
}
