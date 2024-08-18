package main

import (
	"encoding/json"
	"fmt"
)

type Result25 struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func main() {
	var res Result25
	res.Code = 200
	res.Message = "success"
	toJson25(&res)

	setData25(&res)
	toJson25(&res)
}

func setData25(res *Result25) {
	res.Code = 500
	res.Message = "fail"
}

func toJson25(res *Result25) {
	jsons, errs := json.Marshal(res)
	if errs != nil {
		fmt.Println("json marshal error:", errs)
	}
	fmt.Println("json data :", string(jsons))
}
