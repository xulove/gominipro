package app

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"xu.com/xiaowenti/dbs"
)

type App struct {
	//主存储数据库，目前使用pg
	pgdb *sql.DB
	//redis数据库
	rdb redis.Conn
}


func NewApp() (*App,error) {
	pgdb := dbs.GetDB()
	rdb := dbs.GetRDB()
	return &App{pgdb,rdb},nil
}

//设置app处理的mux的一些逻辑
func (app *App) HandlerForMux (mux *http.ServeMux){
	//首页
	mux.HandleFunc("/",app.IndexHandler)
	// 图片上传
	mux.HandleFunc("/uploadimage",app.UploadImageHandler)
	// 图片预览
	mux.HandleFunc("/imagelook/",app.ImageLookHandler)
	//登录
	mux.HandleFunc("/login",app.OnLoginHandler)
	//下单
	mux.HandleFunc("/pay",app.OrderHandler)
	// 在web端进行编辑
	mux.HandleFunc("/web",app.WebIndexHandler)
	//微信支付后的回调
	mux.HandleFunc("/paycallback",app.PayCallbackHandler)
	// 用户所有提问的问题
	mux.HandleFunc("/user/allquestions",app.GetAllQuestionsHandler)
	// 用户所有回答的问题
	mux.HandleFunc("/user/allanswers",app.GetAllAnswerHandler)
}
