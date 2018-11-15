package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"encoding/base64"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister(){
	this.TplName="register.html"
}

func (this *UserController) HandleRegister(){
	userName:=this.GetString("userName")
	password:=this.GetString("password")
	if ""==userName||""==password{
		this.Data["errMsg"]="用户名或密码不能为空！"
		this.TplName="register.html"
		return
	}
	o:=orm.NewOrm()
	var u models.User
	u.UserName=userName
	u.Password=password
	_,err:=o.Insert(&u)
	if err!=nil{
		this.Data["errMsg"]="注册失败，请重新注册！"
		this.TplName="register.html"
		return
	}
	this.Redirect("/login",302)

}

func (this *UserController) ShowLogin(){
	userName:=this.Ctx.GetCookie("userName")
	dec,_:=base64.StdEncoding.DecodeString(userName)
	if string(dec)!=""{
		this.Data["userName"]=string(dec)
		this.Data["checked"]="checked"
	}else {
		this.Data["checked"]=""

	}
	this.TplName="login.html"
}

func (this *UserController) HandleLogin(){
	userName:=this.GetString("userName")
	password:=this.GetString("password")
	if ""==userName||""==password{
		this.Data["errMsg"]="用户名或秘密不能为空！"
		this.TplName="login.html"
		return
	}
	o:=orm.NewOrm()
	var u models.User
	u.UserName=userName
	err:=o.Read(&u,"UserName")
	if err!=nil{
		this.Data["errMsg"]="用户名不存在！"
		this.TplName="login.html"
		return
	}
	if password!=u.Password{
		this.Data["errMsg"]="密码错误！"
		this.TplName="login.html"
		return
	}

	remember:=this.GetString("remember")
	if remember=="on"{
		enc:=base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName",enc,3600)
	}else {
		this.Ctx.SetCookie("userName",userName,-1)
	}
	this.SetSession("userName",userName)

	this.Redirect("/article/articleList",302)
}
func (this *UserController)HandleLogout(){
	this.DelSession("userName")
	this.Redirect("/login",302)

}