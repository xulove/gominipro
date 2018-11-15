package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"xu.com/xiaowenti/dbs"
)

//涉及question相关操作的方法，首先处理相关业务逻辑，然后再调用dbs包下面的数据库方法
// GetAllQuestionsHandler:获取用户的所有问题的请求
// token会在请求参数中传过来
func (app *App) GetAllQuestionsHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println("GetAllQuestionsHandler get request")
	token:=r.FormValue("token")
	if token == ""{
		log.Panicf("no token received on GetAllQuestionsHandler")
	}
	// 更具token获取用户的openid
	openid,err := app.GetOpenidByToken(token)
	if err != nil{
		log.Panicf("GetOpenidByToken error:",err)
	}
	// 根据此openid，获取order中属于此用户的问题
	questions,err := app.GetAllQesByOpenid(openid)
	if err != nil{
		log.Panicf("GetAllQesByOpenid error:",err)
	}
	byteRes,err := json.Marshal(questions)
	if err != nil{
		log.Panicf("questions json marshal error:",err)
	}
	w.Header().Set("content-type","application/json")
	w.Write(byteRes)
}
// 根据openid获取用户的所有提问
func (app *App) GetAllQesByOpenid (openid string)([]Order,error){
	rows:=dbs.Query(app.pgdb,"select * from mini_orders where askuseropenid=$1",openid)
	defer rows.Close()
	var orders []Order = make([]Order,10)
	var order Order
	for rows.Next(){
		err := rows.Scan(&order)
		if err != nil {
			fmt.Printf("rows.Scan error:%v\n",err)
			return nil,err
		}
		orders=append(orders,order)
	}
	return orders,nil
}
