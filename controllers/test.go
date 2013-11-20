package controllers

import (
	"BeeGoTo/models/douban"
	"fmt"
	"github.com/astaxie/beego"
)

type TestController struct {
	beego.Controller
}

func (this *TestController) Get() {

	this.TplNames = "OauthDoubanController/get.tpl"

	db := &douban.Douban{Conf: "conf/douban.conf"}
	if ok := db.Auth("we_get", ""); ok != nil {
		if ok2, url := db.AuthUrl(); ok2 != nil {

		} else {
			this.Data["AuthCodeURL"] = url
		}
	}

	_, dbu := db.User("we_get")
	fmt.Println("dbu:", dbu)

	_, dbtl := db.UserTimeLine("we_get", 20, 0)
	fmt.Println("dbtl:", dbtl)

	//_, dbs := db.Shuo("test")
	//fmt.Println("dbs:", dbs)

}
