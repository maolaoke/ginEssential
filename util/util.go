package util

import (
	"os"
	"github.com/spf13/viper"
)


func InitConfig() {
	workDir, _ := os.Getwd() //获取当前目录
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}