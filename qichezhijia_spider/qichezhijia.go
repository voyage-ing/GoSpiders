package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/axgle/mahonia"
)

// SpiderTemplete 是一个简单的爬虫模版，特别是从resq到body应该是固定的模版；
func SpiderTemplete() {
	fmt.Println("starting crawler: 中文测试")
	objurl := "https://www.autohome.com.cn/beijing/"
	resq, err := http.Get(objurl)
	if err != nil {
		fmt.Println("failed to request url: ", err)
	}
	defer resq.Body.Close() // 函数结束时关闭Body
	body, err := ioutil.ReadAll(resq.Body)

	rawbody := string(body)
	html := ConvertToString(rawbody, "gbk", "utf-8")

	fmt.Println(html)
}

// ConvertToString 是一个解决中文乱码问题，其本质类似与：用gbk编码的中文解码却用utf8解码
// 正确做法是，gbk先用gbk解码再用utf-8编码，然后才能被utf-8解码；
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
