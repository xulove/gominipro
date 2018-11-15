package app

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *App) WebIndexHandler(w http.ResponseWriter,r *http.Request){
	fmt.Println("WebIndexHandler receiver request")
	t, _ :=template.ParseFiles("www/layout.html","www/index.html","www/nav.html","www/wang.html")
	t.ExecuteTemplate(w, "layout","hello world")
}
