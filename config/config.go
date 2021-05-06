package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	ServerAddress    string // 服务监听地址与端口 (e.g. :8090 192.168.1.2:8080)
	AppFileStorePath string // 应用文件存储目录
}

var configuration *Configuration

func initConfiguration() *Configuration {
	configuration = &Configuration{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}

	viper.SetDefault("ServerAddress", ":8080")
	viper.SetDefault("AppFileStorePath", "./app-files")

	if err := viper.Unmarshal(configuration); err != nil {
		log.Fatal(err.Error())
	}

	// 打印当前配置
	log.Println("当前监听地址: ", configuration.ServerAddress)
	log.Println("应用文件存储目录: ", configuration.AppFileStorePath)

	return configuration
}

func GetConfiguration() *Configuration {
	if configuration == nil {
		configuration = initConfiguration()
	}

	return configuration
}
