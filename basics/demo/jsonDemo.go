package demo

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"msgs"`
}

func SetDate(r *Result) {
	r.Code = 500
	r.Message = "fail"
}

func ToJson(r *Result) {
	json, err := json.Marshal(r)
	if(err != nil) {
		fmt.Println(err)
	}
	println(string(json))
}

func JsonDemo() {
	var res Result
	res.Code = 200
	res.Message = "success"

	//序列化
	jsons, errs := json.Marshal(res)
	if errs != nil {
		fmt.Println("json marshal error:", errs)
	}
	fmt.Println("json data :", string(jsons))

	//反序列化
	var res2 Result
	errs = json.Unmarshal(jsons, &res2)
	if errs != nil {
		fmt.Println("json unmarshal error:", errs)
	}
	fmt.Println("res2 :", res2)
}
