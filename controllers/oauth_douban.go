package controllers

import (
	"BeeGoTo/models/douban"
	"fmt"
	"github.com/astaxie/beego"
)

type OauthDoubanController struct {
	beego.Controller
}

func (this *OauthDoubanController) Prepare() {

}

func (this *OauthDoubanController) Get() {

	t_code := this.GetString("code")
	db := &douban.Douban{Conf: "conf/douban.conf"}
	if ok := db.Auth("we_get", t_code); ok != nil {
		if ok2, url := db.AuthUrl(); ok2 != nil {

		} else {
			this.Data["AuthCodeURL"] = url
		}
	}

	_, dbu := db.User("we_get")
	fmt.Println("dbu:", dbu)

	this.Data["dbu"] = dbu

	//_, dbtl := db.UserTimeLine("we_get", 20, 0)
	//fmt.Println("dbtl:", dbtl)
	//_, dbs := db.Shuo("test")
	//fmt.Println("dbs:", dbs)
}

func (this *OauthDoubanController) Post() {

}
