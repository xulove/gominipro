package main

import (
	"fmt"
	"net/http"
	"os"
	"xu.com/xiaowenti/app"
)
func main() {

	app,err := app.NewApp()
	if err != nil {
		fmt.Println("load app failed")
		os.Exit(1)
	}
	mux := http.NewServeMux()
	app.HandlerForMux(mux)
	fmt.Println("http server listen 3009 ")
	if err = http.ListenAndServe(":3009",mux);err != nil{
		fmt.Println("listenAndServe error",err)

	}

}
