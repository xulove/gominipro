package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"xu.com/xiaowenti/dbs"
	"xu.com/xiaowenti/utils"
	"xu.com/xiaowenti/weapp"
)

//处理登录，在小程序无token条件下
/**
处理登录，在程序无token的条件下
1.根据code去取微信的openid
2.根据openid查询本地数据库（自己的服务端）有没有用户的信息
3.若有，则生成token，更新本地并返回给小程序端
4.若无，证明是新用户，生成token，组装用户信息保存在数据库，并返回给小程序端
 */
func (app *App) OnLoginHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println("receive request")
	//获取请求的数据
	code:=r.FormValue("code")
	userInfo := r.FormValue("userinfo")
	// 从配置文件获取小程序的信息
	appid := utils.GetConf("wapp","appid")
	appsecret := utils.GetConf("wapp","appsecret")
	//调用weapp包的login获取openid数据
	lresp,err := weapp.Login(appid,appsecret,code)
	if err != nil{
		fmt.Printf("weapp login error:%v",err)
	}
	// 利用uuid库生成token
	token := uuid.NewV4().String()
	//保存session在redis数据库
	fmt.Println(lresp.OpenID,lresp.SessionKey)
	_,err = app.rdb.Do("SET",token,lresp.OpenID+"_"+lresp.SessionKey)
	if err != nil {
		log.Fatalf("set token error:%v")
	}
	//检查数据库中是否有此openid的用户
	var dbOpenId string
	//err = app.pgdb.QueryRow(`select * from mini_user where openid = $1`,lresp.OpenID).Scan(&dbOpenId)
	err = dbs.QueryRow(app.pgdb,`select openid from mini_user where openid = $1`,lresp.OpenID).Scan(&dbOpenId)
	// 如果有错误，并且错误是sql.ErrNoRows，就说明是是新用户,需要往数据库中插入数据
	if err != nil {
		if err == sql.ErrNoRows{
			fmt.Println("pg have no row about this openid")
			//数据库中没有此openId，
			_,err := dbs.DBExec(app.pgdb,"insert into mini_user(openid,userinfo,token)values($1,$2,$3)",lresp.OpenID,userInfo,token)
			if err != nil{
				log.Panic("插入用户信息失败",err)
			}
		}else{
			log.Panic("select * from mini_user where openid = $1:",err)
		}

	}else{
		fmt.Println("pg have row about this openid")
		// 是老用户，只是重新登录而已。那我们只是更新token就可以了
		dbs.DBExec(app.pgdb,"update mini_user set token = $1 where openid = $2",token,lresp.OpenID)
	}

	//将token返回给小程序端
	byteLresp,err := json.Marshal(struct {
		Token string `json:"token"`
		Error error `json:error`
	}{token,nil})
	if err != nil{
		fmt.Printf("json.Marshal error:%v",err)
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(byteLresp)

}

/**
	获取用户信息，一般是小程序端给个token就可以用获取
 */
func (app *App) OnGetUserinfo(w http.ResponseWriter,r *http.Request){
	fmt.Println("OnGetUserinfo receive request")
}

/**
根据小程序端传来的token在数据库中获取openid
 */
func (app *App) getOpenidByToken1(w http.ResponseWriter,r *http.Request){
	//获取请求的数据
	var err error
	token:=r.FormValue("token")
	var dbtoken string
	err = dbs.QueryRow(app.pgdb,"select token from mini_user where code = $1",token).Scan(&dbtoken)
	var byteResp []byte
	if err != nil {
		byteResp,err = json.Marshal(struct {
			OpenId string `json:openid`
			Error error `json:error`
		}{dbtoken,err})
		if err != nil{
			log.Panic("Json Marshal error :",err)
		}
	}else{

		byteResp,err = json.Marshal(struct {
			OpenId string `json:openid`
			Error error `json:error`
		}{dbtoken,nil})
		if err != nil{
			log.Panic("Json Marshal error :",err)
		}
	}
	w.Header().Set("content-type","application/json")
	w.Write(byteResp)

}

func (app *App) GetOpenidByToken(token string)(string,error){
	var openid string
	err := dbs.QueryRow(app.pgdb,"select openid from mini_user where token = $1",token).Scan(&openid)
	if err != nil{
		fmt.Printf("GetOpenidByToken error:",err)
		return "",err
	}
	return openid,nil
}