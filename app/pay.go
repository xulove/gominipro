package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"xu.com/xiaowenti/dbs"
	"xu.com/xiaowenti/utils"
	"xu.com/xiaowenti/weapp"
	"xu.com/xiaowenti/weapp/payment"
)


/**
下单接口
psot传入token和订单信息
 */
func (app  *App) OrderHandler (w http.ResponseWriter,r *http.Request){
	fmt.Println("OnPayHandler was called")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("read body fail:", err)
		w.WriteHeader(500)
		return
	}
	//获取上传的数据
	var uploadData struct{
		//code是用来获取openid的
		Code string `json:"code"`
		// 下面三个是用来生成问题订单数据的
		Questiontitle string `json:"questiontitle"`
		Questioncontent string `json:"questioncontent"`
		Questionreward int `json:"questionreward"`
	}
	if err = json.Unmarshal(body, &uploadData); err != nil {
		fmt.Println("json Unmarshal fail:", err)
		w.WriteHeader(500)
		return
	}
	// 根据传上来的code，获取openid
	// 从配置文件获取小程序的信息
	appid := utils.GetConf("wapp","appid")
	appsecret := utils.GetConf("wapp","appsecret")
	//调用weapp包的login获取openid数据
	lresp,err := weapp.Login(appid,appsecret,uploadData.Code)
	if err != nil{
		fmt.Printf("weapp login error:%v",err)
	}
	openID := lresp.OpenID
	// 利用uuid库生成商户唯一订单号，这个是不符合规则的，舍弃
	//OutTradeNo := uuid.NewV4().String()
	OutTradeNo := "wx"+utils.RandNewStr(30)
	fmt.Printf("openID:%s,OutTradeNo:%s\n",openID,OutTradeNo)
	// 将订单保存在数据库
	order:= Order{
		openID, uploadData.Questiontitle,uploadData.Questioncontent,uploadData.Questionreward,
		0,
		time.Now().Unix(),
		OutTradeNo,
	}
	fmt.Println("order:",order)
	err = app.OrderSave(order)
	if err != nil{
		log.Panicf("insert order error:",err)
	}
	// 组装支付订单，利用weapp
	// 新建支付订单
	form := payment.Order{
		// 必填
		AppID:      utils.GetConf("wapp","appid"),//"APPID",
		MchID:      utils.GetConf("pay","mch_id"),//"商户号",
		//Body:       "求助>"+order.questiontitle,//"商品描述",
		Body:       "test",//"商品描述",
		NotifyURL:  utils.GetConf("paycall","paycallbackurl"),//"通知地址",
		OpenID:     openID,//"通知用户的 openid",
		OutTradeNo: OutTradeNo,//"商户订单号",
		TotalFee:   order.questionreward*100,//"总金额(分)",

		// 选填 ...
		//IP:        "发起支付终端IP",
		//NoCredit:  "是否允许使用信用卡",
		//StartedAt: "交易起始时间",
		//ExpiredAt: "交易结束时间",
		//Tag:       "订单优惠标记",
		//Detail:    "商品详情",
		//Attach:    "附加数据",
	}
	// 调用微信支付，生成签名
	fmt.Println(form)
	res, err := form.Unify(utils.GetConf("pay","mch_key"))
	if err != nil {
		log.Panicf("统一下单出错：",err)
	}
	fmt.Printf("type of res:%T\n",res)
	fmt.Printf("返回结果: %#v", res)
	// 获取小程序前端调用支付接口所需参数
	params, err := payment.GetParams(res.AppID, utils.GetConf("pay","mch_key"), res.NonceStr, res.PrePayID)
	if err != nil {
		// handle error
		return
	}
	fmt.Println(params)
	byteParams,err := json.Marshal(params)
	if err != nil {
		log.Panicf("json marshal 出错：",err)
	}
	w.Header().Set("content-type","application/json")
	w.Write(byteParams)

}


func (app *App) OnPay(w http.ResponseWriter,r *http.Request){

}
//微信支付后的回调
func (app *App)PayCallbackHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println("PayCallbackHandler receiver request")
}

//将订单保存到数据库
func (app *App)OrderSave(order Order)error{
	fmt.Println(fmt.Sprintf(`insert into mini_orders(askuseropenid,questiontitle,questioncontent,questionreward ,state,createtime,out_trade_no) 
values(%v,%v,%v,%v,%v,%v,%v)`,order.askuseropenid,order.questiontitle,order.questioncontent,order.questionreward ,order.state,order.createtime,order.out_trade_no))
	_,err := dbs.DBExec(app.pgdb,
		`insert into mini_orders(askuseropenid,questiontitle,questioncontent,questionreward ,state,createtime,out_trade_no) 
values($1,$2,$3,$4,$5,$6,$7)`,
		order.askuseropenid,order.questiontitle,order.questioncontent,order.questionreward ,order.state,order.createtime,order.out_trade_no)
	if err != nil{
		fmt.Println("insert order error:",err)
		return err
	}
	return nil

}