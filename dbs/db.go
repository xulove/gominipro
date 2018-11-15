package dbs

import (
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	"log"
	"xu.com/xiaowenti/utils"
)

var DBLog *log.Logger
// 返回pg数据库
func GetDB() (*sql.DB) {
	//从配置文件中加载数据参数
	dbname := utils.GetConf("db","dbname")
	dbuser := utils.GetConf("db","dbuser")
	dbpassword := utils.GetConf("db","dbpassword")
	dbhost := utils.GetConf("db","dbhost")
	//建立数据库对象（注意这只是一个顶层的抽象，真正的连接会在需要时才创建）
	//connStr := "postgres://" + dbuser + ":" + dbpassword + "@" + dbhost + "/" + dbname +"?sslmode=disable"
	connstr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",dbuser,dbpassword,dbname,dbhost)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("getDB error: %v",err)
	}
	//检查数据库是否实际可用
	if err := db.Ping();err != nil{
		log.Fatalf("数据库不是实际可用：%v",err)
	}
	return db
}
// 返回redis数据库
func GetRDB()redis.Conn{
	c,err:=redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("connect redis error:%v",err)
		log.Fatalf("connect redis error:%v",err)
	}
	return c
}

func init() {
	DBLog = utils.GetLog("db")
}
