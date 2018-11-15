package app

// 自己往数据库中保存的数据
type Order struct{
	askuseropenid,questiontitle,questioncontent string
	questionreward int
	state int
	createtime int64
	out_trade_no string
}
// 一个回答的数据结构
type Answer struct {
	ID int `json:"id"`
	AnswerUserOpenid string `json:"answer_user_openid"`
	AnswerContent string `json:"answer_content"`
	QuestionID string `json:"question_id"`
	QuestionTitle string `json:"question_title"`
	CreateTime int64 `json:"create_time"`
}