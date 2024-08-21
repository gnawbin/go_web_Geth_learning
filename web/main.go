package main

import (
	"k8s-web/initiallize"
)

// 项目启动入口
func main() {
	r := initiallize.Routers()
	//initiallize.Viper()
	//initiallize.K8SWithDiscovery()
	//initiallize.InitHarborClient()
	panic(r.Run(":8082"))
}
