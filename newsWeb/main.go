package main

import (
	_ "newsWeb/routers"
	"github.com/astaxie/beego"
	_ "newsWeb/models"
)

func main() {
	beego.AddFuncMap("PrePage",PrePageIndex)
	beego.AddFuncMap("NextPage",NextPageIndex)
	beego.Run()
}
func PrePageIndex(pageIndex int)int{
	prePage:=pageIndex-1
	if prePage<1{
		prePage=1
	}
	return  prePage
}
func NextPageIndex(pageIndex int,pageCount float64)int{
	nextPage:=pageIndex+1
	if nextPage>int(pageCount){
		return  pageIndex
	}
	return  nextPage
}