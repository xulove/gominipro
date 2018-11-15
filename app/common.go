package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"xu.com/xiaowenti/utils"
)
//定义全局变量
var UPLOAD_PATH string

//首页请求逻辑
func (app *App) IndexHandler(w http.ResponseWriter,r *http.Request){
	fmt.Fprintf(w,"hello world")
}
//处理图片上传
func (app *App)UploadImageHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println("starting uploadImageHandler")
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile(UPLOAD_PATH+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在upload目录
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	_,err = io.Copy(f, file)
	if err != nil {
		return
	}
	res := &struct {
		ImageUrl string `json:"imageurl"`
		Error error `json:"error"`
	}{
		handler.Filename,
		nil,
	}
	utils.WriteJson(w,res)
}
//查看图片
func (app *App) ImageLookHandler (w http.ResponseWriter, r *http.Request)  {
	fmt.Println("start ImageLookHandler")
	fmt.Printf(r.URL.Scheme,r.URL.User)
	path := strings.Split(r.URL.Path,"/")
	var name string
	if len(path) > 1 {
		name = path[len(path) - 1]
	} else {
		name = "null.png"
	}
	fmt.Printf(name)
	w.Header().Set("Content-Type", "image/jpg")
	w.Header().Set("Content-Disposition",fmt.Sprintf("inline; filename=\"%s\"",name))

	file, err := ioutil.ReadFile(UPLOAD_PATH + name)
	if err != nil {
		fmt.Fprintf(w,"查无此图片")
		return
	}
	w.Write(file)
}



func init(){
	UPLOAD_PATH  =  utils.GetConf("app","upload_path")
}