# Golang学习-汽车之家爬虫
[TOC]
## 开发过程
### 获取原始信息
```go
func main() {
	fmt.Println("start qichezhijia_spider: ")
	fmt.Println("--------------------------")
	objurl := "https://car.autohome.com.cn/price/tag-25-0-1.html"
	htmlUTF8 := ConvertToString(spider(objurl), "gbk", "utf-8")
	// fmt.Println(htmlUTF8)

	xpath := "/html/body/div[1]/script[2]"
	parsedInfo := XpathParseHtml(htmlUTF8, xpath)
}
```

### 解析json数据
#### 第一层：字符串解析为字典；
var CityItems = {"ProvinceItems":[], "CityItems":[]}

#### 第二层：`CityItems`values里的内容map list；
如`infos["CityItems"] = [map[I:110100 N:北京 P:Beijing S:110000] map[I:120100 N:天津 P:Tianjin S:120000] ...]`
```go
cityitem := cityitems.([]interface{})
for index, citydetails := range cityitem{
	fmt.Println(index, citydetails)
}
```
```bash
0 map[I:110100 N:北京 P:Beijing S:110000]
1 map[I:120100 N:天津 P:Tianjin S:120000]
...
```
#### 第三层：解析每一个map
```bash
0 map[I:110100 N:北京 P:Beijing S:110000]
```
id；areacode；行政区划码；cityname；citypinyin；

```go
id := index
areacode := citydetails.(map[strininterface{})["I"]
cityname := citydetails.(map[strininterface{})["N"]
citypinyin := citydetails.(map[strininterface{})["P"]
fmt.Println(id, areacode, cityname, citypinyin)
```

解析结果：
```
0 110100 北京 Beijing
1 120100 天津 Tianjin
2 130100 石家庄 Shijiazhuang
...
```

### 数据存入Mysql
```go
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
...
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
```
terminal输出：
![](https://tva1.sinaimg.cn/large/007S8ZIlly1gj54vml7mwj31pq0b279k.jpg)

数据库中：
![](https://tva1.sinaimg.cn/large/007S8ZIlly1gj54s58fi9j307x05gdg3.jpg)

## 知识点
1. 中文乱码分析：以GBK为例；如果解码方式用gbk是OK的，但用GBK编码的中文却要用UTF8解码，这显然是不正确的；所以一种解决思路就是，先用gbk解码再用utf8加密；

2. Golang中同一个包下的不同go文件里的函数可以直接使用不用相互导入；这个程序也是这样做的；

3. `fmt.Sprintf` 格式化字符串，并不是打印出来；

4. xpath路径解析html，xml，json：https://github.com/antchfx/xpath；具体针对HTML的使用：https://github.com/antchfx/htmlquery；

5. golang解析json，简单使用可以；解析json列表可以参照：https://studygolang.com/topics/2288；比较复杂并且追求性能："jsoniter"：http://jsoniter.com/index.cn.html

6. mysql链接池，在高并发中特别重要；类似于存储二级缓存的概念；在并发特别高的情况下直接链接mysql server的话，会有安全隐患；这里没涉及这个；

7.  Go import _包名：（如：import _ hello/imp）的作用：当导入一个包时，该包下的文件里所有init()函数都会被执行；有些时候我们并不需要把整个包都导入进来，仅仅是是希望它执行init()函数而已。这个时候就可以使用 import _ 引用该包。即使用【import _ 包路径】只是引用该包，仅仅是为了调用init()函数，所以无法通过包名来调用包中的其他函数。