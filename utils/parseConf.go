package utils

import (
	"gopkg.in/ini.v1"
	"log"
)

const CONFIGNAME  = "app.conf"
//就是为了解析配置文件
func ParseConf(filename string) (*ini.File ){
	cfg, err := ini.Load("app.conf")
	if err != nil {
		log.Fatalf("load app.conf failed:%v",err)
		return nil
	}
	return cfg
}
func GetConf(Section string,key string)string{
	cfg := ParseConf(CONFIGNAME)
	return cfg.Section(Section).Key(key).String()
}
