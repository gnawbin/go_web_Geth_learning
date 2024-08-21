package initiallize

import (
	"k8s-web/global"
	"k8s-web/plugins/harbor"
)

// @Author: morris
func InitHarborClient() {
	enable := global.CONF.System.Harbor.Enable
	scheme := global.CONF.System.Harbor.Scheme
	host := global.CONF.System.Harbor.Host
	username := global.CONF.System.Harbor.Username
	password := global.CONF.System.Harbor.Password
	initHarborClient, err := harbor.InitHarbor(scheme, host, username, password)
	if err != nil && enable {
		panic(err)
	}
	global.HarborClient = initHarborClient
}
