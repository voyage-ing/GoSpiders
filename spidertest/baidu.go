package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("Hello world")
	resq, err := http.Get("http://www.baidu.com/")
	if err != nil {
		fmt.Println("http get err:", err)
		return
	}
	body, err := ioutil.ReadAll(resq.Body)
	if err != nil {
		fmt.Println("read error", err)
		return
	}
	fmt.Println(string(body))
}
