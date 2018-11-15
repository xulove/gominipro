package utils

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	//codes   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/~!@#$%^&*()_="
	// 下面这个是符合微信规则的
	codes   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-|*"
	codeLen = len(codes)
)
func GetLog(prefix string) *log.Logger {
	logFile,err := os.OpenFile("logs/error.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("打开日志文件失败",err)
	}

	newLog := log.New(logFile, prefix, log.Ldate|log.Lshortfile)
	return newLog
}
//往一个http.ResponseWriter里面写入json数据
func WriteJson(w http.ResponseWriter,thing interface{}){
	w.Header().Set("Content-Type","application/json")
	//w.Header().Set("Content-Type","application/json;charset=utf-8")
	respByte,err := json.Marshal(thing)
	if err != nil{
		log.Panic("json.Marshal error:",err)
	}
	w.Write(respByte)
}
// 检测error
func CheckErr(err error){
	if err != nil{
		log.Fatalln(err)
	}
}
// 生成指定长度的随机字符串


func RandNewStr(len int) string {
	data := make([]byte, len)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len; i++ {
		idx := rand.Intn(codeLen)
		data[i] = byte(codes[idx])
	}

	return string(data)
}