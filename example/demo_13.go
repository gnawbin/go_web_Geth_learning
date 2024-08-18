package main

import (
	"encoding/json"
	"fmt"
)

type Result2 struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func main() {
	var res2 Result2
	res2.Code = 200
	res2.Message = "success"
	toJson(&res2)

	setData(&res2)
	toJson(&res2)
}

func setData(res *Result2) {
	res.Code = 500
	res.Message = "fail"
}

func toJson(res *Result2) {
	jsons, errs := json.Marshal(res)
	if errs != nil {
		fmt.Println("json marshal error:", errs)
	}
	fmt.Println("json data :", string(jsons))
}
