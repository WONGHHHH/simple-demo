package main

import (
	"fmt"
	"time"
)

func main() {

	time := time.Now().AddDate(0,0,1).Format("20060102")
    fmt.Println(time)
	// 测试冲突
}