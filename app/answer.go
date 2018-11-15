package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"xu.com/xiaowenti/dbs"
)

func (app *App)GetAllAnswerHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println("GetAllAnswerHandler was called")
	token := r.FormValue("token")
	if token == ""{
		log.Panicf("please give the token")
	}
	// 根据token获取用户的openid
	openid,err := app.GetOpenidByToken(token)
	if err != nil{
		log.Panicf("GetOpenidByToken error:",err)
	}
	// 根据opneid，找到用户的回答
	answers,err := app.GetAllAnsByOpenid(openid)
	if err != nil{
		log.Panicf("GetAllAnsByOpenid error:",err)
	}
	byteRes,err := json.Marshal(answers)
	if err != nil{
		log.Panicf("answers json marshal error:",err)
	}
	w.Header().Set("content-type","application/json")
	w.Write(byteRes)
}
// 根据openid获取用户的所有提问
func (app *App) GetAllAnsByOpenid (openid string)([]Answer,error){
	rows:=dbs.Query(app.pgdb,"select * from mini_answers where answeruseropenid=$1",openid)
	defer rows.Close()
	var answers []Answer = make([]Answer,10)
	var answer Answer
	for rows.Next(){
		err := rows.Scan(&answer)
		if err != nil {
			fmt.Printf("rows.Scan error:%v\n",err)
			return nil,err
		}
		answers=append(answers,answer)
	}
	return answers,nil
}