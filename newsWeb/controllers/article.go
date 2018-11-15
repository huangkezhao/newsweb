package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"path"
	"time"
	"math"
	"strconv"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowArticleList() {
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}

	this.Data["userName"]=userName.(string)
	o := orm.NewOrm()
	var articles []models.Article
	qs := o.QueryTable("Article")
	typeName:=this.GetString("select")
	var count int64
	if typeName==""||typeName=="请选择"{
		count, _ = qs.Count()
	}else{
		count, _ = qs.Filter("ArticleType__TypeName",typeName).Count()

	}
	pageSize := int64(2)

	pageCount := float64(count) / float64(pageSize)
	pageCount = math.Ceil(pageCount)
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount

	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	this.Data["pageIndex"] = pageIndex
	start := pageSize * (int64(pageIndex) - 1)

	if typeName==""||typeName=="请选择"{
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)

	}else{
		this.Data["typeName"] = typeName
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)

	}

	this.Data["articles"] = articles

	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["articleTypes"] = articleTypes
	this.Layout="layout.html"
	this.TplName = "index.html"
}

func (this *ArticleController) ShowAddArticle() {
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["articleTypes"] = articleTypes
	this.Layout="layout.html"

	this.TplName = "add.html"
}

func (this *ArticleController) HandleAddArticle() {
	title := this.GetString("articleName")
	content := this.GetString("content")
	typeId ,_:= this.GetInt("typeName")
	beego.Info("typeName",typeId)
	if "" == title ||0 == typeId || "" == content {
		this.Data["errMsg"] = "文章标题或内容不能为空！"
		this.TplName = "add.html"
		return
	}
	file, head, err := this.GetFile("uploadname")
	defer file.Close()

	if err != nil {
		this.Data["errMsg"] = "获取文件失败！"
		beego.Error(err)
		this.TplName = "add.html"
		return
	}
	if head.Size > 500000 {
		this.Data["errMsg"] = "文件过大！"
		beego.Error(err)
		this.TplName = "add.html"
		return
	}
	fileExt := path.Ext(head.Filename)
	beego.Info(fileExt)
	if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg" {
		this.Data["errMsg"] = "文件格式错误！"
		this.TplName = "add.html"
		return
	}
	fileName := time.Now().Format("2006-01-02-15-04-05") + fileExt

	this.SaveToFile("uploadname", "./static/image/"+fileName)




	o := orm.NewOrm()
	var article models.Article
	var articleType models.ArticleType

	articleType.Id=typeId
	err=o.Read(&articleType)
	if err != nil {
		this.Data["errMsg"] = "获取文章分类失败！"
		beego.Error(err)
		this.TplName = "add.html"
		return
	}

	article.Title = title
	article.Content = content
	article.ArticleType=&articleType
	article.Image = "/static/image/" + fileName
	_, err = o.Insert(&article)
	if err != nil {
		this.Data["errMsg"] = "添加失败！"
		beego.Error(err)
		this.TplName = "add.html"
		return
	}
	this.Redirect("/article/articleList", 302)

}
func (this *ArticleController) ShowArticleDetail() {
	id, err := this.GetInt("id")
	if err != nil {
		errMsg := "请求路径错误！"
		beego.Error(errMsg, err)
		this.Redirect("/article/articleList?errMsg="+errMsg, 302)
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		errMsg := "获取文章失败！"
		beego.Error(errMsg, err)
		this.Redirect("/article/articleList?errMsg="+errMsg, 302)
		return
	}
	var articleType models.ArticleType
	typeId:=article.ArticleType.Id
	beego.Info("typeId",typeId)
	articleType.Id=typeId
	o.Read(&articleType)
	article.ArticleType=&articleType

	m2m:=o.QueryM2M(&article,"Users")
	var user models.User
	userName:=this.GetSession("userName")
	user.UserName=userName.(string)
	o.Read(&user,"userName")
	m2m.Add(user)
	o.LoadRelated(&article,"Users")
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)
	this.Data["users"]=users
	this.Data["article"] = article
	this.Layout="layout.html"

	this.TplName = "content.html"
}
func (this *ArticleController) ShowUpdateArticle() {
	id, err := this.GetInt("id")
	if err != nil {
		errMsg := "请求路径错误！"
		beego.Error(errMsg, err)
		this.Redirect("/article/articleList?errMsg="+errMsg, 302)
		return
	}

	errMsg := this.GetString("errmsg")
	if errMsg != "" {
		this.Data["errMsg"] = errMsg
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		errMsg := "获取文章失败！"
		beego.Error(errMsg, err)
		this.Redirect("/article/articleList?errMsg="+errMsg, 302)
		return
	}
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"]=articleTypes
	this.Data["article"] = article
	this.Layout="layout.html"

	this.TplName = "update.html"

}
func UploadFile(this *ArticleController, tplName string) string {
	file, head, err := this.GetFile("uploadname")
	beego.Info(1233)
	defer file.Close()
	if err != nil {
		this.Data["errMsg"] = "获取文件失败！"
		beego.Error(err)
		this.TplName = tplName
		return ""
	}

	if head.Size > 500000 {
		this.Data["errMsg"] = "文件过大！"
		beego.Error(err)
		this.TplName = tplName
		return ""
	}
	beego.Info(6555)

	fileExt := path.Ext(head.Filename)
	beego.Info(fileExt)
	if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg" {
		this.Data["errMsg"] = "文件格式错误！"
		this.TplName = tplName
		return ""
	}
	fileName := time.Now().Format("2006-01-02-15-04-05") + fileExt

	this.SaveToFile("uploadname", "./static/image/"+fileName)
	return "./static/image/" + fileName
}

func (this *ArticleController) HandleUpdateArticle() {
	id, err := this.GetInt("id")
	title := this.GetString("title")
	content := this.GetString("content")
	filePath := UploadFile(this, "updateArticle")
	beego.Info("title", title)
	beego.Info("content", content)
	beego.Info("filePath", filePath)
	beego.Info(err)
	if title == "" || content == "" || filePath == "" || err != nil {
		errMsg := "文件保存失败！"
		beego.Error(errMsg, err)
		this.Redirect("/article/updateArticle?id="+strconv.Itoa(id)+"&errmsg="+errMsg, 302)
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)

	if err != nil {
		errMsg := "获取文章失败！"
		beego.Error(errMsg, err)
		this.Redirect("/article/updateArticle?id="+strconv.Itoa(id)+"&errmsg="+errMsg, 302)
		return
	}

	article.Content = content
	article.Title = title
	article.Image = filePath
	o.Update(&article)
	this.Redirect("/article/articleList", 302)

}



func (this *ArticleController) HandleDeleteArticle() {
	id, err := this.GetInt("id")
	if err != nil {
		errMsg := "删除错误！"
		beego.Error(errMsg, err)
		this.Redirect("/article/articleList?errMsg="+errMsg, 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	_, err = o.Delete(&article, "Id")
	if err != nil {
		errMsg := "删除失败！"
		beego.Error(errMsg, err)
		this.Redirect("/article/articleList?errMsg="+errMsg, 302)
		return
	}
	this.Redirect("/article/articleList", 302)

}
func (this *ArticleController) ShowAddType() {
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["articleTypes"] = articleTypes
	this.Layout="layout.html"

	this.TplName = "addType.html"

}
func (this *ArticleController) HandleAddType() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		errMsg := "类型名不能为空！"
		beego.Error(errMsg)
		this.Data["errMsg"] = errMsg
		this.Redirect("/article/addType?errMsg="+errMsg, 302)
		return
	}

	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	_, err := o.Insert(&articleType)
	if err != nil {
		errMsg := "添加分类失败！"
		beego.Error(errMsg, err)
		this.Redirect("/article/addType?errMsg="+errMsg, 302)
		return
	}
	this.Redirect("/article/addType", 302)

}
func (this *ArticleController) HandleDeleteType() {
	id, err := this.GetInt("id")
	if err != nil {
		errMsg := "删除失败！"
		beego.Error(errMsg, err)
		this.Redirect("/article/addType?errMsg="+errMsg, 302)
		return
	}
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	_, err = o.Delete(&articleType)
	if err != nil {
		errMsg := "删除失败！"
		beego.Error(errMsg, err)
		this.Redirect("/article/addType?errMsg="+errMsg, 302)
		return
	}
	this.Redirect("/article/addType", 302)
}
