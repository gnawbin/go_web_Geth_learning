package initiallize

import (
	"github.com/spf13/viper"
	"k8s-web/global"
)

func Viper() {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
	err = v.Unmarshal(&global.CONF)
	if err != nil {
		panic(err.Error())
	}
}
