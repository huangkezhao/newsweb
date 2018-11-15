package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"time"
)


type User struct{
	Id int
	UserName string  `orm:"unique"`
	Password string
	Articles []*Article `orm:"rel(m2m)"`

}

type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm："size(100)"`
	Content string`orm:"size(500)"`
	Time time.Time `orm:"type(datetime);auto_now"`
	ReadCount int `orm:"default(0)"`
	Image string `orm:"null"`
	ArticleType *ArticleType `orm:"rel(fk);on_delete(set_null);null"`
	Users []*User `orm:"reverse(many)"`


}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(100)"`
	Articles  []*Article `orm:"reverse(many)"`
}
func init(){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/newsweb?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}