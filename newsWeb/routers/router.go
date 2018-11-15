package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeExec,funcFilter)
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/article/articleList",&controllers.ArticleController{},"get:ShowArticleList")
	beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/article/articleDetail",&controllers.ArticleController{},"get:ShowArticleDetail")
	beego.Router("/article/updateArticle",&controllers.ArticleController{},"get:ShowUpdateArticle;post:HandleUpdateArticle")
	beego.Router("/article/deleteArticle",&controllers.ArticleController{},"get:HandleDeleteArticle")
	beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
	beego.Router("/article/deleteType",&controllers.ArticleController{},"get:HandleDeleteType")
	beego.Router("/article/logout",&controllers.UserController{},"get:HandleLogout")



}
var funcFilter=func(ctx *context.Context){
	userName:=ctx.Input.Session("userName")
	if userName==nil{
		ctx.Redirect(302,"/login")
	}
}